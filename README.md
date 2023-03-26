# Simple File Server in Go and Docker

# How to run this locally
## Without docker
Run `go run main.go`

## With docker, without docker-compose
Run `docker build --tag go-my-fs .`
Create docker volume.
Then, run the image with any port and the volume.

## With docker, with docker-compose
Run `docker-compose build`
Then, run `docker-compose up`
