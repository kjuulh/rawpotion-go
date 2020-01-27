FROM golang:latest
WORKDIR /go/src/github.com/kjuulh/rawpotion-go/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8082
COPY --from=0 /go/src/github.com/kjuulh/rawpotion-go/dist/linux_amd64/release/rawpotion-go .
CMD ["./rawpotion-go"]