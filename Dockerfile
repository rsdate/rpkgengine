FROM golang:1.24.0 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOARCH=arm64 GOOS=darwin go build -ldflags="-X main.version=$(git describe --always --long --dirty)" -o /app/rpkengine

FROM scratch
COPY --from=builder /app/rpkengine /app/rpkengine
ENTRYPOINT ["/app/rpkengine"]
