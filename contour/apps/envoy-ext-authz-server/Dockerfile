FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /envoy-ext-authz-server .

FROM scratch
COPY --from=builder /envoy-ext-authz-server /envoy-ext-authz-server
ENTRYPOINT ["/envoy-ext-authz-server"]
