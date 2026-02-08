FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o csvagg ./cmd/csvagg

FROM alpine:3.21

COPY --from=builder /build/csvagg /usr/local/bin/csvagg

ENTRYPOINT ["csvagg"]
