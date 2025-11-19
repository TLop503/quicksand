#!/bin/bash

# Quicksand - Browser Sandbox Compile & Run Script
set -e  # Exit on any error

echo "🔧 Quicksand - Building and Starting Server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Clean up previous build
echo "🧹 Cleaning previous build..."
rm -f quicksand-app

# Build the application
echo "🏗️  Building Go application..."
go build -o quicksand-app main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Starting Quicksand server on http://localhost:8080"
    echo "📋 Available endpoints:"
    echo "   • http://localhost:8080/ (Web Interface)"
    echo "   • http://localhost:8080/api/health (Health Check)"
    echo "   • http://localhost:8080/api/start (Start Container)"
    echo "   • http://localhost:8080/api/stop (Stop Container)"
    echo "   • http://localhost:8080/api/swap (Swap between Tor and Firefox)"
    echo "   • http://localhost:8080/api/restart (Restart Container)"
    echo ""
    echo "Press Ctrl+C to stop the server"
    echo "=========================================="
    
    # Run the application
    ./quicksand-app
else
    echo "❌ Build failed!"
    exit 1
fi
