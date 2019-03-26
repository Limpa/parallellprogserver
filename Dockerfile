FROM golang:latest

WORKDIR /home/pi/inda/parallellprogserver

COPY . .

RUN go get "github.com/gorilla/websocket"
RUN go install ./...

EXPOSE 8080

CMD [app]
