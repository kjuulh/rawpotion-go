FROM golang:latest
WORKDIR /go/src/github.com/kjuulh/rawpotion-go/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make test