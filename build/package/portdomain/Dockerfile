ARG go_version=1.14.4

FROM golang:$go_version AS builder
WORKDIR /build

# Fetch dependencies
COPY go.mod ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build ./cmd/portdomain

FROM alpine
COPY --from=builder /build/portdomain /bin/portdomain
COPY --from=builder /build/migrations /migrations
COPY --from=builder /build/cmd/portdomain/config.yaml /config/config.yaml

EXPOSE 9000
CMD ["/bin/portdomain"]