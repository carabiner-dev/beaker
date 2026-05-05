# npm runner

The npm runner executes `npm test` in a project directory and parses the
test framework's output to populate a `test-result` in-toto attestation.

It is selected automatically by `beaker run` when a `package.json` is
detected at the root of the project.

## What it runs

```
npm test -- --reporter=tap
```

The output is teed to the user's terminal (so the run looks like a normal
`npm test` invocation) and simultaneously captured for parsing.

## Requirements

For beaker to extract per-test names and pass/fail status, the project's
test command must emit [TAP](https://testanything.org/) on stdout. The
runner appends `-- --reporter=tap` to `npm test`, so the following must
hold:

1. **`package.json` defines a `test` script.** Without it, `npm test`
   exits with an error and no attestation is produced.

2. **The `test` script forwards extra arguments to the underlying test
   framework.** This is the default for plain entries like:

   ```json
   "scripts": { "test": "mocha" }
   ```

   It will not work for hardcoded scripts such as
   `"test": "mocha --reporter=spec"` — the appended `--reporter=tap`
   does not override the earlier flag. In that case, edit the script
   to use TAP, or omit conflicting reporter flags.

3. **The underlying test framework supports a TAP reporter.** Known
   compatible setups:

   | Framework  | How to enable TAP                                  |
   | ---------- | -------------------------------------------------- |
   | Mocha      | `--reporter=tap` (built in)                        |
   | tape       | TAP is the native output                           |
   | node:test  | `node --test --test-reporter=tap` (script must be set up that way) |
   | Jest       | Requires `jest-reporter-tap` (or similar) plugin    |
   | Vitest     | Requires a TAP reporter plugin                      |

   If your framework needs a different flag than `--reporter=tap`, set
   the test script in `package.json` to invoke it directly and ensure
   TAP ends up on stdout.

## Fallback behavior

If no TAP lines are found in the output, the attestation will contain
empty `passedTests` / `failedTests` arrays. The overall `result` field
still reflects the process exit status: `pass` if `npm test` exited
zero, `fail` otherwise.

## Output

The runner produces a `test-result` predicate
(`https://in-toto.io/attestation/test-result/v0.1`) containing:

- `passedTests`: names of TAP test points that reported `ok`
- `failedTests`: names of TAP test points that reported `not ok`
- `result`: `pass` or `fail`
- `configuration`: repository metadata (added by the launcher)

Pass `-a` / `--attest` to wrap the predicate in a full in-toto Statement.
