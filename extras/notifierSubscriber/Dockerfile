FROM alpine:3.12 AS build

RUN apk update \
    && apk add --no-cache go make git gcc musl-dev build-base

WORKDIR /home/app

COPY go.mod go.sum notifiersubscriber.go Makefile /home/app/

RUN make

FROM alpine:3.12 AS prod

RUN mkdir -p /home/app 

WORKDIR /home/app

COPY --from=build /home/app/notifier-subscriber /home/app/

CMD ./notifier-subscriber