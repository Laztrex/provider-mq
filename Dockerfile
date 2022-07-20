FROM golang:latest

ENV APP_NAME controller_mq

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get ./...
RUN go build -ldflags="-s -w" -o ${APP_NAME} cmd/${APP_NAME}/main.go

CMD ./${APP_NAME}