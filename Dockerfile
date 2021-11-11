FROM golang:1.17 as builder

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build .

FROM golang:1.17

WORKDIR /app
COPY --from=builder /app/sillyGirl .
COPY --from=builder /app/develop .
RUN chmod +x /app/sillyGirl

ENTRYPOINT ["/app/sillyGirl"]