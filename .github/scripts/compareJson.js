const fs = require('fs');

// Function to read and parse JSON file
function readJsonFile(filePath) {
    const data = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(data);
}

// Function to compare resources
function compareResources(resource1, resource2) {
    return resource1.type === resource2.type && JSON.stringify(resource1.values) === JSON.stringify(resource2.values) && JSON.stringify(resource1.sensitive_values) === JSON.stringify(resource2.sensitive_values);
}

// Function to compare JSON files
function compareJsonFiles(file1Path, file2Path) {
    const file1 = readJsonFile(file1Path);
    const file2 = readJsonFile(file2Path);

    const resources1 = file1.values.root_module.resources;
    const resources2 = file2.values.root_module.resources;

    let allMatched = true;

    resources1.forEach(resource1 => {
        const match = resources2.find(resource2 => compareResources(resource1, resource2));
        if (match) {
            console.log(`Matching entry found for resource: ${resource1.address}`);
        } else {
            console.log(`No matching entry found for resource: ${resource1.address}`);
            allMatched = false;
        }
    });
}

// Check if the script is run directly from the command line
if (require.main === module) {
    const [,, file1Path, file2Path] = process.argv;

    if (!file1Path || !file2Path) {
        console.error('Please provide paths to two JSON files.');
        process.exit(1);
    }

    const allMatched = compareJsonFiles(file1Path, file2Path);

    console.log(`Value of allMatched: ${allMatched}`);

    if (!allMatched) {
        process.exit(1);
    }
}
