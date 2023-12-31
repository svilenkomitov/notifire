FROM golang:1.20-alpine

COPY . /app
WORKDIR /app

RUN go build -o /notifire cmd/main.go

CMD /notifire