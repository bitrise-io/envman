## Changes

* Environments got new field: IsTemplate. If IsTemplate is true, environment value, will handled with built in go template solutions.
  Exmaple:
  ```
  - script:
    title: Template example
    inputs:
    - content: |-
      {{if .IsCI}}
      echo "CI mode"
      {{else}}
      echo "not CI mode"
      {{end}}
  ```
* Improved environment value and options cast. Now you can use "NO", "Yes", true, false as bool too.


## Install

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/envman/releases/download/1.0.0/envman-$(uname -s)-$(uname -m) > /usr/local/bin/envman
```

Then:

```
chmod +x /usr/local/bin/envman
```
