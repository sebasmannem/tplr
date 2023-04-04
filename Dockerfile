FROM golang:alpine AS tplrbuilder
WORKDIR /usr/src/app

COPY . .
RUN apk update && \
    apk upgrade && \
    apk add git && \
    sh set_version.sh && \
    go mod tidy && \
    go build -o ./bin/tplr ./cmd/tplr

FROM alpine:latest

COPY --from=tplrbuilder /usr/src/app/bin/tplr /usr/local/bin/
ENTRYPOINT [ "/usr/local/bin/tplr" ]
