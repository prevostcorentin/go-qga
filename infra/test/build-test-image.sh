#!/bin/sh

set -euo pipefail

IMAGE_NAME="qga-test"

ctr=$(buildah from debian:stable-slim)

buildah run $ctr -- apt-get update
buildah run $ctr -- apt-get install -y qemu-guest-agent socat
buildah run $ctr -- apt-get clean
buildah run $ctr -- rm -rf /var/lib/apt/lists/*

buildah run $ctr -- wget https://dl.google.com/go/go1.24.5.linux-amd64.tar.gz
buildah run $ctr -- tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz
buildah run $ctr -- rm go1.24.5.linux-amd64.tar.gz

buildah config --env PATH=/usr/local/go/bin:$PATH $ctr
buildah config --env GOPATH=/go $ctr

buildah run $ctr -- mkdir -p /var/run/qemu-ga
buildah run $ctr -- chmod 777 /var/run/qemu-ga

buildah config \
    --entrypoint '["/usr/sbin/qemu-ga","--method=virtio-serial","--path=/var/run/qemu-ga/qga.sock","--logfile=/dev/stderr","--verbose"]' \
    $ctr

buildah commit $ctr $IMAGE_NAME

echo "Image $IMAGE_NAME successfully built"
