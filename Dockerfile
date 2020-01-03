FROM golang:1.13.5 AS builder
WORKDIR /go/src/github.com/jeethsuresh/superman 
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o bin/superman ./cmd/superman

FROM scratch 
COPY --from=builder /go/src/github.com/jeethsuresh/superman/bin/superman .
COPY --from=builder /go/src/github.com/jeethsuresh/superman/data /data
EXPOSE 8080 8080
CMD ["./superman"]