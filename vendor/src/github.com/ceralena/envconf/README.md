envconf
=======

[![Build Status](https://travis-ci.org/ceralena/envconf.svg?branch=master)](https://travis-ci.org/ceralena/envconf)

[![codecov](https://codecov.io/gh/ceralena/envconf/branch/master/graph/badge.svg)](https://codecov.io/gh/ceralena/envconf)

See [godoc](http://godoc.org/github.com/ceralena/envconf) for package documentation.

`envconf` is a [Go](http://golang.org) package that makes it easy to build
explicitly typed, structured configuration objects without complex parsing of
config files or command-line flags.

`envconf` is designed to pull configuration out of the process environment, but
it can be given any function of the type `func(string) string` as an accessor
to the raw config.

Deploy configurable applications without config files; use a shell script or an
init file to set the appropriate config variables and `envconf` will do the
rest of the work at runtime.

An example:

```go
// This program will look up PORT, BIND and BLACKLIST in the process
// environment. If found, it will parse then and set the values on the
// serverConfig object.
package main

import "log"
import "github.com/ceralena/envconf"

func main() {
	var serverConfig struct {
		Port int    `required:"true"`
		Bind string `default:"0.0.0.0"`
		Blacklist []string
	}
	if err := envconf.ReadConfigEnv(&serverConfig); err != nil {
		log.Fatal(err)
	}
}
```

You can alternatively call `envconf.ReadConfigEnvPrefix()` if you want simple
namespacing of the environment variables:

```go
// This program will look up MYSERVER_PORT, MYSERVER_BIND and
// MYSERVER_BLACKLIST in the process environment.
package main

import "log"
import "github.com/ceralena/envconf"

func main() {
	var serverConfig struct {
		Port int    `required:"true"`
		Bind string `default:"0.0.0.0"`
		Blacklist []string
	}
	if err := envconf.ReadConfigEnvPrefix("MYSERVER_", &serverConfig); err != nil {
		log.Fatal(err)
	}
}
```

License
-------

MIT license; see LICENSE.txt.
