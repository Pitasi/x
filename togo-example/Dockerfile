FROM golang:1.24.1-alpine3.21 AS builder
RUN go install anto.pt/x/togo/cmd/togo@latest

FROM alpine:3.21
COPY --from=builder /go/bin/togo /usr/bin/togo
COPY ./content /content
EXPOSE 8080
CMD ["togo", "/content", "-laddr", "8080"]
