FROM golang:1.18-alpine as builder
WORKDIR /app
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg
COPY go.* ./
RUN go mod download
# RUN go get -d -v ./...
# RUN go vet -v
# RUN go test -v
RUN CGO_ENABLED=0 go build -o app ./cmd/kubesnap
RUN apk add -U curl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl

FROM alpine:latest
RUN apk add -U --no-cache git ansible
COPY --from=builder /app/app ./
COPY --from=builder /app/kubectl /usr/local/bin/kubectl
CMD ["./app"]