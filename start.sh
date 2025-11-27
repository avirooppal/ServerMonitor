#!/bin/bash
set -e

echo "Building Frontend..."
cd web
npm install
npm run build
cd ..

echo "Copying Frontend assets..."
rm -rf cmd/server/dist
mkdir -p cmd/server/dist
cp -r web/dist/* cmd/server/dist/

echo "Building Backend..."
go build -o server-moni ./cmd/server

echo "Starting Server..."
./server-moni
