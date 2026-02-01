#!/bin/sh
set -e

echo "Building docker-secrets container image..."

# Build the image
docker build -t docker-secrets:latest .

echo ""
echo "Build complete!"
echo ""
echo "Image: docker-secrets:latest"
echo ""
echo "To run:"
echo "  docker run --rm docker-secrets:latest --help"
echo ""
echo "To check image size:"
echo "  docker images docker-secrets:latest"
