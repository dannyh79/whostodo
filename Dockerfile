# syntax=docker/dockerfile:1

FROM golang:1.22.1 AS dev
WORKDIR /app
VOLUME /app

FROM dev AS build
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /app
COPY . .
RUN go mod download \
  && go build -o /whostodo

FROM gcr.io/distroless/base-debian12:nonroot AS build-release
WORKDIR /
COPY --from=build /whostodo /whostodo
EXPOSE 8080

ENTRYPOINT ["./whostodo"]
