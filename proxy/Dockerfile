FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /proxy-server

EXPOSE 5000

CMD [ "/proxy-server" ]