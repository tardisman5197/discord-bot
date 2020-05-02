FROM golang:1.14

RUN mkdir -p /go/src/github.com/tardisman5197/discord-bot
WORKDIR  /go/src/github.com/tardisman5197/discord-bot

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o discord-bot

ENTRYPOINT ["./discord-bot"]