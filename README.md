# Madlib Service

This small webservice includes a single endpoint (`/madlib`) which returns a
sentence template randomly filled with words pulled from
`https://reminiscent-steady-albertosaurus.glitch.me/`.

## Install

1. Clone this repo.

```
$ git clone git@github.com:purple4reina/madlib.git
```

2. Build the docker image.

```
$ docker build -t madlib
```

3. Run the service with docker.

```
$ docker run -p 8080:8080 madlib
```

4. Test the service to make sure it is working

```
$ curl http://localhost:8080/madlib
```

## Development

If you would like to contribute, you can run the tests with

```
$ go test ./...
```

Running the service locally can be done with

```
$ go run madlib.go
```
