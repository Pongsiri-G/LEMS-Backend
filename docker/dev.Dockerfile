FROM golang:1.25-alpine

# Create a non-root user and group
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

RUN apk add make

# Create app directory with proper ownership
RUN mkdir -p /app && chown -R appuser:appuser /app

# Copy with read-only permissions (chmod in COPY command)
COPY --chown=root:root --chmod=555 . /app/

# Switch to non-root user
USER appuser

WORKDIR /app

# Create writable cache directories for Go (not in /app)
RUN mkdir -p /home/appuser/go /home/appuser/.cache

# Install air
RUN go install github.com/air-verse/air@latest

ENV PATH="/home/appuser/go/bin:${PATH}"
ENV GOCACHE="/home/appuser/.cache"

CMD ["air", "-c", "/app/.air.toml"]