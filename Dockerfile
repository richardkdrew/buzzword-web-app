# DOCKER VERSION 1.12.0
FROM golang:1.6.2-alpine

MAINTAINER Richard Drew <richardkdrew@gmail.com>

# set the start/working directory
WORKDIR /app

# copy the files into place
COPY . /app

# expose the hosting port
EXPOSE 8080/tcp

# set the run command
ENTRYPOINT ["./main"]
