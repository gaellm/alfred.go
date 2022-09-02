<img align="right" width="500px" src="https://repository-images.githubusercontent.com/447978404/5ba210c2-9e56-463a-88f2-7aa1a8ea55b2">

# Alfred.go
Because even software engineer super heros needs a good valet. Afred.go is a mock, written in Go (Golang), for performance testing. Alfred manages a mock list, offers helpers, permits to trigger asynchronous actions, and offers the ability to wrap users' javascript functions; users have infinite creatives possibilities.


The main goal is to provide a simple way to mock APIs (json, xml, plain text), without developpement knowledges, with a minimum resources footprint _(thinked for cloud projects)_. Designed for performance, Alfred.go provides metrics, and tracing _(using [OpenTelemetry.io](https://opentelemetry.io))_ for observability.

# Quick reference
- Maintained by:  
[The Alfred.go project](https://github.com/gaellm/alfred.go)
- Where to get some documentation:  
[Alfred.go's documentation](https://gaellm.github.io/alfred.go/)

# How to use this image.

This image contains Alfred.go with the default configuration and some examples. It uses a [distroless base image](https://github.com/GoogleContainerTools/distroless), so it contains only the application and its runtime dependencies. The image not contains package managers, shells or any other programs you would expect to find in a standard Linux distribution. 

## Pull latest image

```console
$ docker pull gaellm/alfred.go
```

## Start a container

```console
$ docker run --rm --name alfred.go -p 8080:8080 gaellm/alfred.go
```
Then access http://localhost:8080/ to display handled mocks. You can find the [defaults provided mocks here](https://github.com/gaellm/alfred.go/tree/main/user-files/mocks) with corresponding requests examples. And the [default configuration file here](https://github.com/gaellm/alfred.go/blob/main/configs/config.json).

## My first mock
The easiest way to manage mocks is to create a folder, and put all mocks json files in it. 
```console
$ mkdir my-mocks
$ echo '{"request":{"url":"/whereareyou"},"response":{"body":"Gotham, Sir"}}' > my-mocks/city.json
```
Then mount this folder as volume:
```console
$ docker run --rm --name alfred.go \
  -p 8080:8080 \
  -v $PWD/my-mocks/:/alfred/user-files/mocks/ \
  gaellm/alfred.go

```

Access http://localhost:8080/whereareyou to use the mock. Default mocks, functions and body-files folders path can be set using Alfred's configuration.

## Configuration

### From file

Alfred.go's image configuration file path is _/alfred/configs/config.json_, so you can use volumes to overwrite the [default file](https://github.com/gaellm/alfred.go/blob/main/configs/config.json).
```console
$ echo '{"alfred":{"log-level":"DEBUG"}}' > my-config.json

$ docker run --rm --name alfred.go \
  -p 8080:8080 \
  -v $PWD/my-config.json:/alfred/configs/config.json \
  gaellm/alfred.go

```
### From environment

However each option can be set via environment variables. So, to set log level to debug as seen above, just set _ALFRED_LOG_LEVEL=debug_:
```console
$ docker run --rm --name alfred.go \
  -p 8080:8080 \
  -e ALFRED_LOG_LEVEL=debug \
  gaellm/alfred.go

```

# License

View [license information](https://raw.githubusercontent.com/gaellm/alfred.go/main/LICENSE) for the software contained in this image.

As with all Docker images, these likely also contain other software which may be under other licenses (such the base image, along with any direct or indirect dependencies of the primary software being contained).

As for any pre-built image usage, it is the image user's responsibility to ensure that any use of this image complies with any relevant licenses for all software contained within.
