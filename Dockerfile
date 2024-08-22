FROM golang:alpine AS build_base

RUN apk update --no-cache && apk add git
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_base as build

WORKDIR /app

RUN apk update --no-cache && apk add git
ADD ./ /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o sample-app  .


FROM alpine
WORKDIR /app
COPY --from=build /app/sample-app /app

ENTRYPOINT ["/app/sample-app"]

EXPOSE 8080