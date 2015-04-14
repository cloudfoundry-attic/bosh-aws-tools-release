# bosh-aws-tools

When running the tests, they will use `bosh` in whatever state it's already
configured for.

## Running the tests

You will need a working Go environment with `$GOPATH` set, and you will need
`bosh` in your `$PATH`.

See [Go][go] for instructions on installing `go`.

### Configuration

Before running the tests, you must make sure you've targetted the desired Director
and can create a dev release:

```
bosh target <host>
bosh create release
```

You must also set `$CONFIG` to point to a `.json` file which contains the
configuration for the tests.

For example:

```sh
cat > config.json <<EOF
{
  "aws_access_id": "access",
  "aws_secret_access_key": "secret",
  "route53_zone_names": ["bosh.domain.com."]
}
EOF
```

### Running

```sh
export CONFIG=$PWD/config.json
./bin/test [ginkgo arguments ...]
```

The `test` script will pass any given arguments to [ginkgo](https://github.com/onsi/ginkgo), so this is where
you pass `-focus=`, `-nodes=`, etc.

#### Running in parallel

To run the tests in parallel, pass `-nodes=X`, where X is how many examples to
run at once.

```sh
./bin/test -nodes=10
```

#### Seeing command-line output

If you want to see the output of all of the commands it shells out to, set
`VERBOSE_OUTPUT` to `true`.

```sh
export VERBOSE_OUTPUT=true
./bin/test
```

[go]: http://golang.org
