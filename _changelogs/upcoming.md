## Changes

* Removed unnecessary log :'[ENVMAN] - Failed to execute command: XYZ'
* Run command exist code fix: if executing the command failed, but command exit code was 0, envman run exit code was 0, instead of valid exit code.
* envman configs file added to $HOME/.envman/configs.json path. Configs contains the maximum allowed environment value size in KB (20 KB default) and the maximum allowed environment list size in KB (100 KB default).
* Environment variables became new field `skip_if_empty`, if this property is set and the value is empty, the Environment will not exported during `bitrise run`. *This is the same behavior as previous version of bitrise worked*. To add env with skip_if_empty, just call `envman add --key KEY --value VALUE --skip-if-empty`.
