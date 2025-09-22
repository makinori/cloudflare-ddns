FROM golang:1.25.1 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
go build -o /app/cloudflare-ddns .

# ---

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt \
/etc/ssl/certs/ca-certificates.crt

COPY --from=build /app/cloudflare-ddns /cloudflare-ddns

CMD ["/cloudflare-ddns"]
