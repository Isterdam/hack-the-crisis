FROM golang:1.14

COPY . /go/src/github.com/hack-the-crisis/backend/
WORKDIR /go/src/github.com/hack-the-crisis/backend/src
COPY . .

RUN go get -d -v ./...
COPY . .
EXPOSE 8080

CMD ["go", "run", "main.go"]