FROM golang:1.24 AS buildgo
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o ./build/hostamat ./cmd/hostamat

FROM alpine:3
WORKDIR /app
COPY --from=buildgo /app/build/hostamat /app/hostamat
EXPOSE 1001
ENTRYPOINT ["/app/hostamat"]
