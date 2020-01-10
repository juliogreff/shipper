#!/bin/bash -ex

docker tag $DOCKER_REGISTRY/shipper:$IMAGE_TAG bookingcom/shipper:$IMAGE_TAG
docker tag $DOCKER_REGISTRY/shipper:$IMAGE_TAG bookingcom/shipper:latest
docker push bookingcom/shipper

docker tag $DOCKER_REGISTRY/shipper-state-metrics:$IMAGE_TAG bookingcom/shipper-state-metrics:$IMAGE_TAG
docker tag $DOCKER_REGISTRY/shipper-state-metrics:$IMAGE_TAG bookingcom/shipper-state-metrics:latest
docker push bookingcom/shipper-state-metrics
