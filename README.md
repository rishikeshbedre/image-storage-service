# IMAGE STORAGE SERVICE

![Licence](https://img.shields.io/github/license/rishikeshbedre/image-storage-service)
[![Build Status](https://travis-ci.com/rishikeshbedre/image-storage-service.svg?branch=master)](https://travis-ci.com/rishikeshbedre/image-storage-service)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/129c09fa009440928ba88410be8d5fd1)](https://app.codacy.com/manual/rishikeshbedre/image-storage-service?utm_source=github.com&utm_medium=referral&utm_content=rishikeshbedre/image-storage-service&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/rishikeshbedre/image-storage-service)](https://goreportcard.com/report/github.com/rishikeshbedre/image-storage-service)

Image Storage Service is a microservice based on REST APIs to store and retrieve images. It features REST end-points to add/delete/modify albums and images. It is written using [Gin Web Framework](https://github.com/gin-gonic/gin) and [jsoniter](https://github.com/json-iterator/go) to make server high performant and have used [Swaggo](https://github.com/swaggo/swag) for API documentation.

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

<p align="justify">`*But when ran simple tests, I noticed that while adding a big image the container's RAM usage was increasing and not releasing it back. So I debug the application using pprof and saw that both heap and CPU profile were fine. Later I decreased the value of memory block used by multipart forms (from 32mb to 8mb), this decreased the memory usage but didn't solve the memory issue. As I was using golang version above 1.12, this is the expected behaviour (memory is freed by GC but OS doesn't take it back until its required). This can be harmful in container environment, so I have set 'GODEBUG=madvdontneed=1' instead of MADV_FREE and calling FreeOSMemory() in a interval.*`</p>

## Docker

## Testing

## Limitations