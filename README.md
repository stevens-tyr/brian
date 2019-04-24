# brian 🧠

Brian the Binary is the Brain that runs concurrent goroutines to exec files in the Kubernetes job based on the test data provided in the env.

**THINGS TO KNOW**:

- Tests are _not_ going to run in order. This is on purpose, so different tests can run concurrently.
- Brian gets angry if a test takes longer than 2 minutes and will panic.

## TODO

- ~Add Exec to Run Build Script~
- ~unzip/untar supporting/submission files to `/tmp/job`~
- ~Report Grades back to `plague_doctor`~
- ~Swap panics with errors~
- support zip files

## Example env

```
BACKEND_URL="http://localhost:5000/api/v1/plague_doctor" # plague doctor endpoint
COURT_HERALD_URL="http://localhost:4000/api/v1/grader"
SUB_ID="5cc0d5b6d1f3acfda73f5cb4"
ASSIGN_ID="5cc0c4e3823d153ae0e3201d"
BUILD_CMD="echo -n foobar"
TESTS="<Stringified JSON Below>"
JOB_SECRET="foobar" # shared with plague doctor and court herald
```

## Example Tests JSON Input for Brian:

```json
[
  {
    "name": "Testing Echo",
    "expectedOutput": "foobar",
    "studentFacing": true,
    "testCMD": "echo -n foobar"
  },
  {
    "name": "Testing Echo Multiline",
    "expectedOutput": "foo\nbar",
    "studentFacing": true,
    "testCMD": "echo -n foo\nbaz"
  }
]
```

## Example JSON Output for Brian:

```json
[
  {
    "id": 0,
    "name": "Test 2",
    "panicked": false,
    "passed": true,
    "output": "foobar",
    "html": "<span>foobar</span>",
    "testCMD": "echo -n foobar",
    "name": "Testing Echo"
  },
  {
    "id": 1,
    "name": "Test 2",
    "panicked": false,
    "passed": false,
    "output": "foobar",
    "html": "<span>foo</span><del style=\"background:#ffe6e6;\">&para;<br></del><span>ba</span><del style=\"background:#ffe6e6;\">z</del><ins style=\"background:#e6ffe6;\">r</ins>",
    "testCMD": "echo -n foobar",
    "name": "Testing Echo Again"
  }
]
```

If a string does not match, the html will have a pretty diff like so:

<pre style="fontFamily: monospace; background-color: white; color: black;">
<span>foo</span><del style="background-color:#ffe6e6;">&para;<br></del><span>ba</span><del style="background-color:#ffe6e6;">z</del><ins style="background-color:#e6ffe6;">r</ins>
</pre>
