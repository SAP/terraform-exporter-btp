const assert = require('node:assert/strict');
const { test } = require('node:test');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const { compareStates } = require('./compareJson');

const identity = { subaccount_id: 'account-1', app_name: 'app', plan_name: 'standard' };
const instance = (id = identity) => ({ identity: id, identity_schema_version: 0 });
const resource = (id = identity) => ({
    mode: 'managed', type: 'btp_subaccount_subscription', name: 'subscription_0',
    instances: [instance(id)],
});
const state = (...resources) => ({ version: 4, resources });
const compare = (reference, actual) => compareStates(reference, actual, () => {});

test('matches identities despite key order, resource order, names, attributes and state metadata', () => {
    const reference = state(resource(), resource({ subaccount_id: 'account-2' }));
    reference.serial = 1;
    reference.lineage = 'old-lineage';
    reference.resources[0].instances[0].attributes = { id: 'old-id', last_modified: 'yesterday' };
    const actual = state(resource({ subaccount_id: 'account-2' }),
        resource({ plan_name: 'standard', app_name: 'app', subaccount_id: 'account-1' }));
    actual.serial = 99;
    actual.lineage = 'new-lineage';
    actual.resources[1].name = 'subscription_42';
    actual.resources[1].instances[0].attributes = { id: 'new-id', last_modified: 'today' };
    actual.resources[1].instances[0].identity_schema_version = 1;
    assert.equal(compare(reference, actual), true);
});

test('compares every field of a composite identity, even when attributes are equal', () => {
    for (const field of Object.keys(identity)) {
        assert.equal(compare(state(resource()), state(resource({ ...identity, [field]: 'different' }))), false);
    }
    assert.equal(compare(state(resource()), state(resource({ ...identity, extra: 'field' }))), false);
});

test('requires matching resource type and mode', () => {
    for (const change of [{ type: 'other_type' }, { mode: 'data' }]) {
        assert.equal(compare(state(resource()), state({ ...resource(), ...change })), false);
    }
});

test('includes module resources and all count/for_each instances regardless of address', () => {
    const grouped = { ...resource(), module: 'module.parent.module.child', instances: [
        { ...instance({ id: 'one' }), index_key: 0 },
        { ...instance({ id: 'two' }), index_key: 'key' },
    ] };
    const split = state(resource({ id: 'two' }), { ...resource({ id: 'one' }), name: 'renamed' });
    assert.equal(compare(state(grouped), split), true);
    assert.equal(compare(state(grouped), state(resource({ id: 'one' }))), false);
});

test('rejects missing, unexpected and surplus duplicate instances in either direction', () => {
    for (const [left, right] of [
        [state(resource()), state()],
        [state(resource()), state(resource(), resource({ id: 'extra' }))],
        [state(resource(), resource()), state(resource())],
    ]) {
        assert.equal(compare(left, right), false);
        assert.equal(compare(right, left), false);
    }
    assert.equal(compare(state(resource(), resource()), state(resource(), resource())), true);
    assert.equal(compare(state(), state()), true);
});

test('rejects absent, empty and all-null identities instead of falling back to attributes', () => {
    for (const id of [undefined, null, {}, [], 'id', { id: null }, { scope: { id: null } }]) {
        const invalid = state({ ...resource(), instances: [{ identity: id, attributes: { id: 'same' } }] });
        assert.throws(() => compare(invalid, state(resource())), /Missing or empty identity/);
        assert.throws(() => compare(state(resource()), invalid), /Missing or empty identity/);
    }
});

test('allows null optional identity fields alongside populated fields', () => {
    const optional = state(resource({ id: 'one', optional_name: null }));
    assert.equal(compare(optional, structuredClone(optional)), true);
    assert.equal(compare(optional, state(resource({ id: 'two', optional_name: null }))), false);
});

test('rejects show JSON and malformed raw state', () => {
    for (const invalid of [null, {}, { version: 3, resources: [] }, { version: 4 },
        { format_version: '1.0', values: { root_module: { resources: [] } } }]) {
        assert.throws(() => compare(invalid, state()), /terraform state pull/);
    }
    assert.throws(() => compare(state({ ...resource(), instances: null }), state()), /Invalid resource/);
    assert.throws(() => compare(state({ ...resource(), instances: [null] }), state()), /Invalid instance/);
});

test('rejects deposed instances from unfinished replacements', () => {
    const pending = state({ ...resource(), instances: [{ ...instance(), deposed: 'replacement' }] });
    assert.throws(() => compare(pending, pending), /Deposed instance/);
});

test('reports resource addresses without logging identity or attribute values', () => {
    const messages = [];
    const reference = state(resource({ id: 'private-reference-id' }));
    const actual = state({ ...resource({ id: 'private-actual-id' }), name: 'unexpected' });
    assert.equal(compareStates(reference, actual, message => messages.push(message)), false);
    assert.match(messages.join('\n'), /No matching identity.*subscription_0/);
    assert.match(messages.join('\n'), /Unexpected resource.*unexpected/);
    assert.doesNotMatch(messages.join('\n'), /private-/);
});

test('CLI returns success or failure and handles invalid input without exposing state', t => {
    const directory = fs.mkdtempSync(path.join(os.tmpdir(), 'compare-state-'));
    t.after(() => fs.rmSync(directory, { recursive: true, force: true }));
    const referencePath = path.join(directory, 'reference.json');
    const actualPath = path.join(directory, 'actual.json');
    const run = (...args) => spawnSync(process.execPath, [path.join(__dirname, 'compareJson.js'), ...args], { encoding: 'utf8' });
    fs.writeFileSync(referencePath, JSON.stringify(state(resource())));
    fs.writeFileSync(actualPath, JSON.stringify(state(resource())));
    let result = run(referencePath, actualPath);
    assert.equal(result.status, 0);
    assert.match(result.stdout, /All resource identities matched/);
    fs.writeFileSync(actualPath, JSON.stringify(state()));
    assert.equal(run(referencePath, actualPath).status, 1);
    fs.writeFileSync(actualPath, '{"secret":"do-not-log-this"');
    result = run(referencePath, actualPath);
    assert.equal(result.status, 1);
    assert.match(result.stderr, /Invalid JSON/);
    assert.doesNotMatch(result.stderr, /do-not-log-this/);
    assert.equal(run(referencePath, path.join(directory, 'missing.json')).status, 1);
    assert.equal(run().status, 1);
});
