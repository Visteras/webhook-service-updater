FROM golang:1.9 as builder

RUN mkdir -p $GOPATH/src/github.com/Visteras/webhook-service-updater
COPY ./ $GOPATH/src/github.com/Visteras/webhook-service-updater
WORKDIR $GOPATH/src/github.com/Visteras/webhook-service-updater

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook-service-updater .

FROM alpine

RUN apk update && apk add docker

RUN mkdir -p /app/files/yml

COPY --from=builder /go/src/github.com/Visteras/webhook-service-updater/webhook-service-updater /app/webhook-service-updater
RUN chmod +x /app/webhook-service-updater
EXPOSE 4000
WORKDIR /app

CMD ["/app/webhook-service-updater"]