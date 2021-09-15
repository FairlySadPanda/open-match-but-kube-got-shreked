# Development Guide

This dev guide was written for the original Open Match on Kube, and has been doctored to remove redundant information irrelevant to pulling up Open Match locally for testing.

## Install Prerequisites

To build Open Match you'll need the following applications installed.

 * [Git](https://git-scm.com/downloads)
 * [Go](https://golang.org/doc/install)
 * Make (Mac: install [XCode](https://itunes.apple.com/us/app/xcode/id497799835))
 * [Docker](https://docs.docker.com/install/) including the
   [post-install steps](https://docs.docker.com/install/linux/linux-postinstall/).

On Debian-based Linux you can install all the required packages (except Go) by
running:

```bash
sudo apt-get update
sudo apt-get install -y -q make google-cloud-sdk git unzip tar
```

*It's recommended that you install Go using their instructions because package
managers tend to lag behind the latest Go releases.*

Note that at the moment the GCP SDK is a requirement, but this should be removable in time.

## Building code and images

```bash
# Reset workspace
make clean
# Run tests
make test
# Build all the images.
make build-images -j$(nproc)
```

_**-j$(nproc)** is a flag to tell make to parallelize the commands based on
the number of CPUs on your machine._