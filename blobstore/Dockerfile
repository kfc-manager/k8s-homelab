FROM --platform=linux/arm64 arm64v8/golang:1.22.5-alpine3.20 as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY main.go ./

RUN go build -o bin main.go

FROM --platform=linux/arm64 arm64v8/alpine:3.20

COPY --from=build /app/bin /main

ENTRYPOINT [ "/main" ]
