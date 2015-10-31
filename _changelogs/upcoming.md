## Changes

* Removed unnecessary log :'[ENVMAN] - Failed to execute command: XYZ'
* Run command exist code fix: if executing the command failed, but command exit code was 0, envman run exit code was 0, instead of valid exit code.
