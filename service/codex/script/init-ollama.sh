#!/bin/bash

echo "Pulling nomic-embed-text model..."
docker exec doris-codex-ollama ollama pull nomic-embed-text
echo "Done!" 