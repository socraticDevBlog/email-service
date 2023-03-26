FROM golang:1.20.2 AS builder

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux

WORKDIR /go/src/app

COPY go.mod ./

RUN go mod download

COPY src/ ./

RUN go build -o main .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/main .

CMD [ "./main" ]