FROM golang:1.25-alpine

# Create a non-root user and group
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

RUN apk add make

# Create app directory with proper ownership
RUN mkdir -p /app && chown -R appuser:appuser /app

COPY --chown=appuser:appuser . /app/

RUN chmod -R 555 /app

# Switch to non-root user
USER appuser

WORKDIR /app

# Create writable cache directories for Go (not in /app)
RUN mkdir -p /home/appuser/go /home/appuser/.cache

# Install air as non-root user (will install to /home/appuser/go/bin)
RUN go install github.com/air-verse/air@latest

# Ensure the air binary is in PATH
ENV PATH="/home/appuser/go/bin:${PATH}"
ENV GOCACHE="/home/appuser/.cache"

CMD ["air", "-c", "/app/.air.toml"]