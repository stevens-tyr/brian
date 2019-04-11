# brian 🧠

Brian the Binary is the Brain that runs concurrent goroutines to exec files in the Kubernetes job based on the test data provided in the env.

**THINGS TO KNOW**:

- Tests are _not_ going to run in order. This is on purpose, so different tests can run concurrently.
- Brian gets angry if a test takes longer than 2 minutes and will panic.

## TODO

- Add Exec to Run Build Script
- unzip/untar supporting/submission files to `/tmp/job`
- Report Grades back to `plague_doctor`
- Swap panics with errors

## Example env

```
API_URI="localhost:5000/api/v1/plague_doctor" # plague doctor endpoint
TEST_DATA="<Stringified JSON of Schema Below>"
JOB_SECRET="foobar" # shared with plague doctor and court herald
```

## Example JSON Input for Brian:

```json
{
  "submissionID": "<MONGO ID>",
  "assignmentID": "<MONGO ID>",
  "testBuildCMD": "echo -n build stuff here",
  "tests": [
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
}
```

## Example JSON Output for Brian:

```json
[
  {
    "panicked": false,
    "passed": true,
    "output": "foobar",
    "html": "<span>foobar</span>",
    "testCMD": "echo -n foobar",
    "name": "Testing Echo"
  },
  {
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
