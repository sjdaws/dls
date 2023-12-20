FROM golang:1.21.3 as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine

COPY --from=builder /app/dls /app/dls
RUN mkdir -p /app/database
COPY database/migrate /app/database/migrate

EXPOSE 80

CMD ["/app/dls"]
