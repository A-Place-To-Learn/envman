# envman

**Technology Preview:** this repository is still under active development, breaking changes are expected and feedback is greatly appreciated!

## Who and for what?

- connect tools with each other : one tool saves an ENV, the other uses it (input & output)
- manage your environment-sets : use different ENVs for different projects
- complex environment values : if you want to store a complex input as ENV (for example a change log text/summary), `envman` makes this easy, you don't have to encode the value so that when you call bash's `export KEY=value` it will store your value as you intended. You can just give `envman` the value as it is through a `--valuefile` option or as an input stream (or just edit the related `.envstore` file manually), no encoding required.
- switch between environment sets : if you work on multiple projects where each one requires different environments you can manage this with `envman`

## Install

Check the latest release for instructions at: [https://github.com/bitrise-io/stepman/releases](https://github.com/bitrise-io/stepman/releases)


## How? - Use cases

- multi PATH handling: you have `packer` in your `$HOME` dir, in a bin subdir and terraform in another
	- create an envman `.envset` to include these in your `$PATH`


## Develop & Test in Docker

* Build: `docker build -t envman .`
* Run: `docker run --rm -it -v `pwd`:/go/src/github.com/bitrise-io/envman --name envman-dev envman /bin/bash`

Or use the included scripts:

* To build&run: `bash _scripts/docker_build_and_run.sh`
* Once it's built you can just run it: `bash _scripts/docker_run.sh`


## How does it work?

`envman` will run the command you specify
with the environments found in it's environment store.

When you add a new key-value pair with `envman add` it stores
the key and value in an `.envstore` file in the current
directory (this can be changed), and when you call `envman run`
the next time the environments in the `.envstore` file will
be loaded for the command, and the command will be executed
with an environment which includes the specified environment
variables.

This is the same as you would manually set all the
environments, one by one with `export KEY=value` (in bash)
before calling the command.

`envman` makes it easy to switch between environment sets
and to isolate these environment sets from each other -
you don't have to `unset` environments, the specified
environment set will only be available for that single
run of `envman` / the command and won't affect subsequent
calls of the command or `envman`.


## Usage example: Simple Bash example

Add/store an environment variable:

```
envman add --key MY_TEST_ENV_KEY --value 'test value for test key'
```

Echo it:

```
envman run bash -c 'echo "Environment test: $MY_TEST_ENV_KEY"'
```

*Why `bash -c` is required? Because `echo` in itself
does not do any environment expansion, it simply prints
it's input. So if you want to see the value of an environment
you have to run it through a tool/shell which actually
performs the environment expansion (replaces the environment
key with the value associated with it).*


## Usage example: Ruby

### Add environment variable with `--value` flag

```
system( "envman add --key SOME_KEY --value 'some value' --expand false" )
```

### Add environment variable from an input stream

```
IO.popen('envman add --key SOME_KEY', 'r+') {|f|
	f.write('some value')
	f.close_write
	f.read
}
```

### Add environment variable with a value file

```
require 'tempfile'

file = Tempfile.new('SOME_FILE_NAME')
file.write('some value')
file.close

system( "envman add --key SOME_KEY --valuefile #{file.path}" )

file.unlink
```

## Usage example: Go

### Add environment variable with `--value` flag

```
import "os/exec"

cmd := "envman add --key SOME_KEY --value 'some value'"
c := exec.Command("bash", "-c", cmd)
err := c.Run()
if err != nil {
   // Handle error
}
```

### Add environment variable from an input stream

```
import "os/exec"

cmd := "echo 'some value' | envman add --key SOME_KEY --expand false"
c := exec.Command("bash", "-c", cmd)
err := c.Run()
if err != nil {
	// Handle error
}
```

### Add environment variable with a value file

```
import (
	"os/exec"
	"fmt"
)

cmd := fmt.Sprintf("envman add --key SOME_KEY --valuefile /path/to/file/which/contains/the/value")
c := exec.Command("bash", "-c", cmd)
err := c.Run()
if err != nil {
	// Handle error
}
```

## Usage example: Bash

### Add environment variable with `--value` flag

```
envman add --key SOME_KEY --value 'some value'
```

### Add environment variable from an input stream

```
echo "some value" | envman add --key SOME_KEY
```

### Add environment variable with a value file

```
envman add --key SOME_KEY --valuefile /path/to/file/which/contains/the/value --expand false
```
