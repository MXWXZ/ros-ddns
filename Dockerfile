FROM golang:latest AS build
WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -o ros-ddns main.go

FROM scratch AS prod
WORKDIR /app
COPY --from=build /build/ros-ddns .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV GIN_MODE=release

EXPOSE 8080
CMD ["./ros-ddns"]
