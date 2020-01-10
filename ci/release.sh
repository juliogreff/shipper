#!/bin/bash -e

docker tag $DOCKER_REGISTRY/shipper:$IMAGE_TAG juliogreff/shipper:$IMAGE_TAG
docker push juliogreff/shipper:$IMAGE_TAG

docker tag $DOCKER_REGISTRY/shipper-state-metrics:$IMAGE_TAG juliogreff/shipper-state-metrics:$IMAGE_TAG
docker push juliogreff/shipper-state-metrics:$IMAGE_TAG
