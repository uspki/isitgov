# isitgov 
is a simple RESTful server that polls the registration information on home.dotgov.gov/data periodically and provides REST interfaces (binding to :8080) to obtain information for all .gov domain registrations, a single registration (using the domain name as a key) or to determine if a domain is owned by the Federal Government or by State/Local Government.

To build this project, make sure you have a current version of Golang installed and configured, `cd $GOPATH`, `go get github.com/uspki/isitgov`, `cd $GOPATH/src/github.com/uspki/isitgov` and finally build the binary with `go get ./...`

If you encounter any dependency errors, manually run `go get github.com/gorilla/mux` and those issues should be resolved.

The three interfaces are available at localhost:8080/registrations (requiring no key and returning all registrations in JSON), localhost:8080/registrations/ (requiring a single domain - not JSON formatted - as a key and returning the registration information for it) and localhost:8080/isStateLocal/ (requiring no key and returning all registrations in JSON), localhost:8080/registrations/ (requiring a single domain - not JSON formatted - as a key and returning `true` if the domain belongs to state or local government or `false` if it either does not or has not been registered).
