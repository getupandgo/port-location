ARG go_version=1.14.4

FROM golang:$go_version AS builder
WORKDIR /build

# Fetch dependencies
COPY go.mod ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build ./cmd/clientapi

FROM alpine
COPY --from=builder /build/clientapi /bin/clientapi
COPY --from=builder /build/port_data /port_data
COPY --from=builder /build/cmd/clientapi/config.yaml /config/config.yaml

EXPOSE 8000
CMD ["/bin/clientapi"]