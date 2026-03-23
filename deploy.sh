#!/usr/bin/env bash
set -e

echo "🚀 Starting deployment..."

# Pull latest code
echo "📦 Pulling latest code..."
git pull origin main

# Build new image
echo "🔨 Building Docker image..."
docker compose -f docker-compose.prod.yml build

# Run migration
echo "🗄️  Running database migration..."
docker compose -f docker-compose.prod.yml run --rm migrate

# Restart API service
echo "♻️  Restarting API service..."
docker compose -f docker-compose.prod.yml up -d api

# Clean up old images
echo "🧹 Cleaning up..."
docker image prune -f

echo "✅ Deployment completed successfully!"
echo "📊 Check logs: docker logs -f lems-api"
