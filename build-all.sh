#!/bin/sh

APP_NAME="raygun"
#VERSION="v0.1.15"

# Build for all platforms
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}_linux_amd64

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o ${APP_NAME}_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o ${APP_NAME}_darwin_arm64

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o ${APP_NAME}_windows_amd64.exe

echo "Build complete!"

