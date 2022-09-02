<a href="https://github.com/gaellm/alfred.go/actions/workflows/build.yml"><img align="right" src="https://github.com/gaellm/alfred.go/actions/workflows/build.yml/badge.svg?branch=main"></a><br>

<img align="right" width="500px" src="https://repository-images.githubusercontent.com/447978404/5ba210c2-9e56-463a-88f2-7aa1a8ea55b2">

# Alfred.go
Because even software engineer super heros needs a good valet. Afred.go is a mock, written in Go (Golang), for performance testing. Alfred manages a mock list, offers helpers, permits to trigger asynchronous actions, and offers the ability to wrap users' javascript functions; users have infinite creatives possibilities.


The main goal is to provide a simple way to mock APIs, without developpement knowledges, with a minimum resources footprint _(thinked for cloud projects)_. Designed for performance, Alfred.go provides metrics, and tracing _(using [OpenTelemetry.io](https://opentelemetry.io))_ for observability.

## Easy as pie

### Install & Start
You can Download a release [here](https://github.com/gaellm/alfred.go/releases), and extract the archive content in a folder to start Alfred _(or use the [Docker image](https://hub.docker.com/r/gaellm/alfred.go))_. No dependencies needed. Ok that's all for install ;) , then execute alfred.go binary. Congratulations you're not alone anymore, Alfred will help you to test your apps!

### My first mock
A mock is a simple json file. Create it in the _mocks_ folder _(default path is user-files/mocks/)_: 

_mine.json:_
```json
{
    "request": {
        "url": "/mine"
    }
}
```
and restart Alfred.go. Now each GET request to http://localhost:8080/mine will answer a '200 OK' status code. For sure you can add response body, headers, or use variables, random fakers, sent callback requests etc ... [The detailed  documentation](https://gaellm.github.io/alfred.go/) will help you to create awsome mocks ;)

## Get Started
Go to the documentation _Get Started_ section [here](https://gaellm.github.io/alfred.go/)!
