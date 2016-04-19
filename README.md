[![Build Status](http://beta.drone.io/api/badges/lgtmco/lgtm/status.svg)](http://beta.drone.io/lgtmco/lgtm)

LGTM is a simple pull request approval system using GitHub protected branches and maintainers files. Pull requests are locked and cannot be merged until the minimum number of approvals are received. Project maintainers can indicate their approval by commenting on the pull request and including LGTM, looks good to me, in their approval text. For more information please see the documentation at https://lgtm.co/docs

### Setup

Please see our [installation guide](https://lgtm.co/docs/install/) to install the official Docker image.

### Build

Clone the repository to your Go workspace:

```sh
git clone git://github.com/lgtmco/lgtm.git $GOPATH/src/github.com/lgtmco/lgtm
cd $GOPATH/src/github.com/lgtmco/lgtm
```

Commands to build from source:

```sh
export GO15VENDOREXPERIMENT=1

make deps    # Download required dependencies
make gen     # Generate code
make build   # Build the binary
```

If you are having trouble building this project please reference its .drone.yml file. Everything you need to know about building LGTM is defined in that file.
