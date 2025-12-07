#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notely

gcloud builds submit --tag us-central1-docker.pkg.dev/notely-479616/notely-ar-repo/notely:latest .
