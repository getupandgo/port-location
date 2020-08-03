# Port location service

### Running using docker-compose:

`task compose_up` will build docker images for both services, spawn database and services

### Running integration tests:

`task run_integration`

### Lint proto files and go code 

`task lint`

### Generate mocks and grpc code 

`task generate`

`task gen_grpc`

### Improvements
 - database normalization (see comments in `internal/portdomain/storage`)
 - graceful service shutdown
 - more convenient integration testing environment (see comments in `internal/portdomain/server/post_test.go`)
 - better error handling and logging
 - unit tests for transport layer for input validation
 - mapping internal error and grpc error
 - more convenient migration usage

### Running raw binaries (just in case):

1) Create config file with local values (`config.local.yaml`) in both `cmd/clientapi` 
and `cmd/portdomain`.

2) Spawn postgreSQL database for `portdomain` service.

3) Build services with `task build_bin`

4) Run corresponding binaries by `task run_portdomain_bin`, `task run_clientapi_bin`. 

Note that ClientApi service depends on portdomain service, and portdomain service on db.
