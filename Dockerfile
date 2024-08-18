FROM golang:1.23

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV TODO_PORT=7540

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE ${TODO_PORT}

RUN  go build -o /my_app

CMD ["/my_app"]