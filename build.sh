#!/bin/sh
set -e

echo "Building secrets-sync container image..."

# Build the image
docker build -t secrets-sync:latest .

echo ""
echo "Build complete!"
echo ""
echo "Image: secrets-sync:latest"
echo ""
echo "To run:"
echo "  docker run --rm secrets-sync:latest --help"
echo ""
echo "To check image size:"
echo "  docker images secrets-sync:latest"
