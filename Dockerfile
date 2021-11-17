###
FROM k8s.gcr.io/pause:3.1 AS pause

###
FROM golang:1.17-alpine3.14 AS builder

WORKDIR $GOPATH/src/github.com/primeroz/kafka-benchmarker
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

### Final Image
FROM alpine:3.14

COPY --from=builder /go/bin/producer /usr/local/bin
COPY --from=builder /go/bin/consumer /usr/local/bin
COPY --from=pause /pause /pause

RUN apk add --no-cache jq bash curl

USER nobody

ENTRYPOINT ["/pause"]



