FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 go build -o csvagg ./cmd/csvagg

FROM alpine:3.21

COPY --from=builder /build/csvagg /usr/local/bin/csvagg

ENTRYPOINT ["csvagg"]
