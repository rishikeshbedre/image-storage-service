# IMAGE STORAGE SERVICE

![Licence](https://img.shields.io/github/license/rishikeshbedre/image-storage-service)
[![Build Status](https://travis-ci.com/rishikeshbedre/image-storage-service.svg?branch=master)](https://travis-ci.com/rishikeshbedre/image-storage-service)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/129c09fa009440928ba88410be8d5fd1)](https://app.codacy.com/manual/rishikeshbedre/image-storage-service?utm_source=github.com&utm_medium=referral&utm_content=rishikeshbedre/image-storage-service&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/rishikeshbedre/image-storage-service)](https://goreportcard.com/report/github.com/rishikeshbedre/image-storage-service)

Image Storage Service is a microservice based on REST APIs to store and retrieve images. It features REST end-points to add/delete/modify albums and images. It is written using [Gin Web Framework](https://github.com/gin-gonic/gin) and [jsoniter](https://github.com/json-iterator/go) **to make server high performant** and have used [Swaggo](https://github.com/swaggo/swag) for API documentation.

## Contents

- [IMAGE STORAGE SERVICE](#image-storage-service)
  - [Contents](#contents)
  - [How it Works](#how-it-works)
  - [Usage](#usage)
  - [API Documentation](#api-documentation)
  - [Performance Metrics](#performance-metrics)
  - [Docker](#docker)
    - [Docker Build](#docker-build)
    - [Docker Run](#docker-run)
    - [Docker Compose](#docker-compose)
  - [Testing](#testing)
  - [Limitations](#limitations)

## How it Works

![design](https://github.com/rishikeshbedre/image-storage-service/blob/master/extras/design.jpg)

<p align="justify">Image storage service has REST end-points to add/delete/modify albums and images. These images are stored in database(file-system) along with meta-info like which image belongs to which album. On successful operations of storing or retrieving of images or albums this service publishes notification messages to MQTT broker and anyone who wants to listen this notification can subscribe to the topic 'imagestore/notifier'.</p>

## Usage

To install Image Storage Service, you need to install [Go](https://golang.org/)(**version 1.12+ is required**) and set your Go workspace.

1. This project uses go modules and provides a make file. You should be able to simply install and start:

```sh
$ git clone https://github.com/rishikeshbedre/image-storage-service.git
$ cd image-storage-service
$ make
$ mkdir image-db
$ export HOSTIP=$(hostname -I | awk '{print $1}')
$ ./image-storage-service
```

2. Prior to starting image storage service, you need to install and start [MQTT broker](https://mosquitto.org/blog/2013/01/mosquitto-debian-repository/)

## API Documentation

Image storage service serves the API document using swagger and this interactive document can be accessed by:

```sh
$ http://host-address:3333/swagger/index.html
```

Or else link to the document is shown in the logs. Swagger document can be disabled by setting 'PRODMODE' env variable.

## Performance Metrics

Benchmarking for this application is not done.

<p align="justify"><i>"As this application uses Gin web framework, the default logs of gin server shows how much time is consumed by each request to send response back. By running sample tests, there was high RAM usage when adding the image to storage and this issue was debugged by running pprof on heap and CPU profile. Later this issue was solved by decreasing the value of memory block used by multipart forms (from 32mb to 8mb). As this microservice will run in container environment and this application is written using golang version above 1.12, 'GODEBUG=madvdontneed=1' flag is set instead of default MADV_FREE and calling FreeOSMemory() in a interval so that memory is freed whenever the application doesn't need it anymore."</i></p>

## Docker

### Docker Build

To build the container for this microservice, you need to run this script:

```sh
$ git clone https://github.com/rishikeshbedre/image-storage-service.git
$ cd image-storage-service
$ ./extras/build.sh
```

This script builds the image storage service and then pulls eclipse-mosquitto image from docker-hub and then builds notification subscriber service image which helps in subscribing to notification.

### Docker Run

After the docker build stage, you can run the containers using following commands:

```sh
$ docker run -it -p 1883:1883 --name mq_broker_service eclipse-mosquitto:1.6.9
$ docker run -it -p 3333:3333 -e HOSTIP=`hostname -I | awk '{print $1}'` -v "$(pwd)"/image-db:/home/app/image-db  --name image-store image-storage-service:0.0.1
$ docker run -it --name notify-sub -e HOSTIP=`hostname -I | awk '{print $1}'`  notification-subscriber-service:0.0.1
```

### Docker Compose

After the docker build stage, you can run the containers using following script:

```sh
$ ./extras/run.sh
```

## Testing

To run test just run following command:

```sh
$ go mod download
$ make test
```

## Limitations

1. There is no authentication on api routes and the server is running on http mode.
2. There is no field validation or image file validation on api routes.
3. As the microservice uses file-system as database, security has to be maintained at OS.
4. Notification service has to be run in secure mode.