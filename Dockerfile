FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY targets.txt targets.txt
COPY cmd/ cmd/
COPY internal/ internal/

RUN go install github.com/g0ldencybersec/gungnir/cmd/gungnir@latest

RUN CGO_ENABLED=0 go build -o cert_inspector ./cmd/cert_inspector

# Runtime stage: use minimal Alpine image with nonroot user
FROM alpine:3.18

# Create user with UID 65532 (same as distroless nonroot)
RUN addgroup -g 65532 nonroot && adduser -D -u 65532 -G nonroot nonroot

WORKDIR /app

# Copy built in files and binaries
COPY --from=build /app/targets.txt /app/targets.txt
COPY --from=build /app/cert_inspector /app/cert_inspector
COPY --from=build /go/bin/gungnir /app/gungnir

# Create logs folder and fix permissions (won't help when volume mounted, but good fallback)
RUN mkdir /app/logs && chown -R nonroot:nonroot /app/logs

USER nonroot

CMD ["./cert_inspector", "-binary", "/app/gungnir", "-targets", "targets.txt", "-log-dir", "/app/logs"]
