#!/usr/bin/env bash

set -o nounset
set -o xtrace
set -o errexit

#
sudo yum update -y
sudo yum install -y docker
docker -v
sudo systemctl start docker
sudo systemctl enable docker
systemctl status docker

exit $?
