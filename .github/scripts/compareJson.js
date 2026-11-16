const fs = require('node:fs');
const _ = require('lodash');

function readJsonFile(filePath) {
    const data = fs.readFileSync(filePath, 'utf8');
    try {
        return JSON.parse(data);
    } catch {
        // Parser errors can include state values, so do not print the original error.
        throw new Error(`Invalid JSON in ${filePath}`);
    }
}

function isObject(value) {
    return value !== null && typeof value === 'object' && !Array.isArray(value);
}

function hasIdentityValue(value) {
    if (value === null || value === undefined) return false;
    if (typeof value === 'object') return Object.values(value).some(hasIdentityValue);
    return true;
}

// Read raw state, including instances in modules and count/for_each resources.
function resourceInstances(state) {
    if (!isObject(state) || state.version !== 4 || !Array.isArray(state.resources)) {
        throw new Error('Expected Terraform raw state version 4 from terraform state pull. Regenerate references created with terraform show -json.');
    }

    return state.resources.flatMap(resource => {
        if (!isObject(resource) || typeof resource.type !== 'string' ||
            typeof resource.name !== 'string' || !['managed', 'data'].includes(resource.mode) ||
            !Array.isArray(resource.instances)) {
            throw new Error('Invalid resource in Terraform state.');
        }
        const address = [resource.module, resource.mode === 'data' ? 'data' : null, resource.type, resource.name]
            .filter(Boolean).join('.');

        return resource.instances.map(instance => {
            if (!isObject(instance)) throw new Error(`Invalid instance for ${address}`);
            const instanceAddress = instance.index_key === undefined
                ? address : `${address}[${JSON.stringify(instance.index_key)}]`;
            if (!isObject(instance.identity) || !hasIdentityValue(instance.identity)) {
                throw new Error(`Missing or empty identity for ${instanceAddress}. Refresh the state with a provider that supports resource identities.`);
            }
            if (instance.deposed) {
                throw new Error(`Deposed instance for ${instanceAddress}. Complete the pending replacement before comparing states.`);
            }
            return {
                address: instanceAddress,
                mode: resource.mode,
                type: resource.type,
                identity: instance.identity,
            };
        });
    });
}

function compareResources(reference, actual) {
    return reference.mode === actual.mode && reference.type === actual.type &&
        _.isEqual(reference.identity, actual.identity);
}

function compareStates(reference, actual, log = console.log) {
    const expected = resourceInstances(reference);
    const remaining = resourceInstances(actual);
    let allMatched = true;

    for (const resource of expected) {
        const index = remaining.findIndex(candidate => compareResources(resource, candidate));
        if (index === -1) {
            log(`No matching identity found for resource: ${resource.address}`);
            allMatched = false;
        } else {
            // Consume each match so a single instance cannot satisfy multiple references.
            remaining.splice(index, 1);
            log(`Matching identity found for resource: ${resource.address}`);
        }
    }

    for (const resource of remaining) {
        log(`Unexpected resource: ${resource.address}`);
        allMatched = false;
    }
    return allMatched;
}

function compareJsonFiles(file1Path, file2Path) {
    return compareStates(readJsonFile(file1Path), readJsonFile(file2Path));
}

if (require.main === module) {
    const [,, file1Path, file2Path] = process.argv;

    if (!file1Path || !file2Path) {
        console.error('Please provide paths to two raw Terraform state JSON files.');
        process.exitCode = 1;
    } else {
        try {
            const allMatched = compareJsonFiles(file1Path, file2Path);
            console.log(allMatched
                ? 'Summary: All resource identities matched after import.'
                : 'Summary: Resource identities did not match after import.');
            if (!allMatched) process.exitCode = 1;
        } catch (error) {
            console.error(`State comparison failed: ${error.message}`);
            process.exitCode = 1;
        }
    }
}

module.exports = { compareStates };
