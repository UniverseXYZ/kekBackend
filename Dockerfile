FROM golang:1.15 AS build

WORKDIR /barnbridge

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM scratch
COPY --from=build /barnbridge/barnbridge-backend /barnbridge-backend
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/barnbridge-backend", "run", "--config=/config/config.yml"]
