FROM golang:1.14 as builder

WORKDIR /go/src/github.com/rezaAmiri123/service-user

COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine:latest

RUN apk --no-cache add ca-certificates
ENV config=docker
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/rezaAmiri123/service-user .

CMD ["./service-user"]
