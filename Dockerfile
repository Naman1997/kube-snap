FROM golang:1.18 as builder
WORKDIR /app
COPY cmd/kube-snap/ ./
RUN go mod download
RUN go build -a -installsuffix cgo -o app .

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/app ./
CMD ["./app"]