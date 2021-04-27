FROM golang:1.16 AS build

WORKDIR /kekDAO

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM scratch
COPY --from=build /kekDAO/abis /abis
COPY --from=build /kekDAO/dashboard/web /dashboard/web
COPY --from=build /kekDAO/kekBackend /kekBackend
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/kekBackend", "run", "--config=/config/config.yml"]
