# Start by building the application.
FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/alfred.go

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
LABEL Name="alfred.go" Version="v1.0-alpha" 

WORKDIR /alfred

COPY --from=build /go/bin/alfred.go .
COPY configs ./configs
COPY user-files ./user-files

USER 1000
EXPOSE 8080

ENTRYPOINT ["./alfred.go"]