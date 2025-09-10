#!/bin/bash

echo "🚀 Starting Logistics API..."

# Load environment variables
export $(cat .env | xargs)

# Build and run
echo "📦 Building application..."
go build -o bin/server cmd/server/main.go

echo "🌟 Starting server..."
./bin/server