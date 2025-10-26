FROM golang:1.25-alpine

# Install necessary packages without caching to keep the image small
RUN apk add --no-cache make ca-certificates tzdata \
	&& update-ca-certificates

# Create an unprivileged user and group to avoid running as root
RUN addgroup -S app \
	&& adduser -S -G app -h /home/app app

# Set working directory
WORKDIR /app

# Copy source code and ensure ownership belongs to the non-root user
COPY --chown=app:app . /app/

# Install Air (hot reloader)
RUN go install github.com/air-verse/air@latest \
	&& chown -R app:app /app /go/bin /go/pkg || true

# Drop privileges
USER app

# Expose the default app port (optional; adjust if different)
# EXPOSE 8080

CMD ["air", "-c", "/app/.air.toml"]
