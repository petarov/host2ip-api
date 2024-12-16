FROM golang:1.23 AS build

WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags "-s" -o host2ip

FROM alpine:latest AS run
RUN apk --update add ca-certificates tzdata

WORKDIR /root
COPY --from=build /app/host2ip .

EXPOSE 7029

CMD ["./host2ip", "-addr", "0.0.0.0", "-apikey", "somerandomstring"]