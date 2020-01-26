FROM golang:alpine
WORKDIR /go/src/github.com/kjuulh/rawpotion/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8082
COPY --from=0 /go/src/github.com/kjuulh/rawpotion/app .
CMD ["./app"]