#!/usr/bin/env sh

echo ">>> Generating the binary ..."
# This the tool used to build the binary in otelcol-dev.
# It generates all the code calling the factories defined by the modules in the configfile
ocb --config builder-config.yaml

echo ">>> Building docker image ..."
DOCKER_BUILDKIT=1 docker build . -t otel-contrib

echo ">>> Pushing image to registry ..."
docker tag otel-contrib eu.gcr.io/halfpipe-io/ee-o11y/otel-collector:latest
docker push eu.gcr.io/halfpipe-io/ee-o11y/otel-collector:latest

docker tag otel-contrib eu.gcr.io/halfpipe-io/ee-o11y/otel-collector:0.90.1-o11y1
docker push eu.gcr.io/halfpipe-io/ee-o11y/otel-collector:0.90.1-o11y1