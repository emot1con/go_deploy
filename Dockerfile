FROM golang:1.24 as bukder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd

FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/main /
CMD ["/main"]