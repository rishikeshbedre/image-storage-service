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
  - [Docker](#docker)
    - [Docker Build](#docker-build)
    - [Docker Run](#docker-run)
    - [Docker Compose](#docker-compose)
  - [Testing](#testing)
  - [Limitations](#limitations)

## How it Works

![design](https://github.com/rishikeshbedre/image-storage-service/blob/master/extras/design.jpg)

