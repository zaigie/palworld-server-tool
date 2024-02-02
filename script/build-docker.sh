#!/bin/bash
echo "Building pst Docker image..."
docker build -f Dockerfile --build-arg proxy=true --build-arg version=$(git describe --tags) -t palworld-server-tool .

echo "Building pst-agent Docker image..."
docker build -f Dockerfile.agent --build --build-arg proxy=true -t palworld-server-tool-agent .