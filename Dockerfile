FROM golang:alpine as builder

ENV APP_NAME provider_mq

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

COPY ./go.mod ./go.sum ./
RUN go mod download

RUN go build -ldflags="-s -w" -o ${APP_NAME} cmd/${APP_NAME}/main.go

CMD []

FROM alpine

ENV APP_NAME provider_mq

WORKDIR /go/src/${APP_NAME}

COPY --from=builder /go/src/${APP_NAME} /go/src/${APP_NAME}

CMD ./${APP_NAME}