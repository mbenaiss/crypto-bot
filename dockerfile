FROM golang:alpine3.11 as build-stage
RUN apk add build-base
WORKDIR /project
COPY . .
RUN GOOS=linux go build -mod vendor -ldflags "-s -w" -o app cmd/main.go

FROM alpine:3.11.3
WORKDIR /project
COPY --from=build-stage /project /project
ENTRYPOINT ["./app"]
