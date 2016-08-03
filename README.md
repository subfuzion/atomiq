# amp

At the top level, we still have the `swarm` script for starting infrastructure services. We will need this for one more sprint, but it will be going away.

## Directories

* api/rpc - contains the directories for specifying our rpc services and messages. They should be segregated by service and only contain *.proto files for definitions and *.pb.go files that are generated by the rpc rule in the Makefile (`make rpc`). Running `make install` or `make build` will automatically invoke the rpc rule.
* api/server - contains the implementation files for the api/rpc/* service interfaces, and also contains `server.go`, which registers each service with a server on port 50051.
* api/cmd - contains directories corresponding to the two generated executable binaries: `amp` (the CLI) and `amplifier` (the server daemon).
* data/elasticsearch - the elasticsearch data layer (obviously).
* vendor - used by glide to lock down versions and is supposed to be committed to version control, along with glide.lock and glide.yaml.

## Tests

Test should generally be colocated with the packages they test. For example, see `api/server/project_test.go`, which tests the project service
(`api/server/project.go`).

It is usually most expedient to just run `go test` in the package directory you want to test. But for convenience, `make test` will run tests
on the `github.com/appcelerator/amp/api/server` package to ensure that service tests are passing.

## Makefile

* `make` - (no arguments) will print version/build info, run a check on your code, then install `amp` and `amplifier` in `$GOPATH/bin`.
* `make check` - runs a series of formatting and lint checks on the source. Make sure to run before submitting a PR.
* `make fmt` - will format the source according to go conventions.
* `make clean` - cleans up auto-generated code and deletes the `amp` and `amplifier` executables from `$GOPATH/src`.
* `make rpc` - run this if you want to update generated files without doing a full build.
* `make install` - will (re)build the project and install the executable binaries (`amp` and `amplifier`) in `$GOPATH/src`.
* `make build` - will build a docker image (`appcelerator/amp:latest`) that contains the binaries.
* `make run` - will run `amplifier` in a container.

## Notes

Use the `swarm` shell script to launch amp cluster services.

    $ ./swarm pull
    $ ./swarm start
    $ ./swarm ls
    $ ./swarm restart (equivalent to stop, pull, start)
    $ ./swarm stop
