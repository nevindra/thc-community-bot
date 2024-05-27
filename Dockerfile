FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY . .

RUN go build -o main .

FROM debian:bookworm-slim

WORKDIR /app

# copy the ca-certificate.crt from the build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /app/main

CMD ["/app/main"]