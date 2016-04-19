LGTM is a simple pull request approval system using GitHub protected branches and simple MAINTAINERS files. For more information please see the documentation https://lgtm.co/docs

## Setup

Please see our [installation guide](https://lgtm.co/docs/install/) to install the official Docker image.

## Build

Clone the repository to your Go workspace:

```sh
git clone git://github.com/drone/drone.git $GOPATH/src/github.com/drone/drone
cd $GOPATH/src/github.com/drone/drone
```

Commands to build from source:

```sh
export GO15VENDOREXPERIMENT=1

make deps    # Download required dependencies
make gen     # Generate code
make build   # Build the binary
```

## Help

Contributions, questions, and comments are welcomed and encouraged. LGTM developers hang out in the [lgtmco/lgtm](https://gitter.im/lgtmco/lgtm) room on gitter. We ask that you please post your questions to gitter before creating an issue.


If you are having trouble building this project please reference its .drone.yml file. Everything you need to know about building LGTM is defined in that file.
