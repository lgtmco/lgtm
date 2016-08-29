[![Build Status](http://beta.drone.io/api/badges/lgtmco/lgtm/status.svg)](http://beta.drone.io/lgtmco/lgtm)

LGTM is a simple pull request approval system using GitHub protected branches and maintainers files. Pull requests are locked and cannot be merged until the minimum number of approvals are received. Project maintainers can indicate their approval by commenting on the pull request and including LGTM (looks good to me) in their approval text. For more information please see the documentation at https://lgtm.co/docs

### Status

LGTM is actively used by thousands of repositories. The lack of commit activity in the Git repository is not an indication that the project is abandoned or inactive. The lack of activity is an indication the author views this project as being feature-complete.

### Development

LGTM is meant to be extremely simple and focused and is largely considered feature-complete. The author is certainly interested in minor improvements and bug fixes, but is not interested in major enhancements. Feel free to fork the project and extend (and even re-brand) as you see fit.

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

### Usage

#### .lgtm file

Each repository managed by LGTM can have a .lgtm file in the root of the repository. This file provides the configuration properties that LGTM uses for the repository. If the .lgtm file is missing then default values are used for all properties.

#### MAINTAINERS file

Each repository managed by LGTM should have a MAINTAINERS file that specifies who is allowed to approve pull requests.

For simple approval, you can use a flat-file format for the list of approvers:
```
Brad Rydzewski <brad.rydzewski@mail.com> (@bradrydzewski)
Matt Norris <matt.norris@mail.com> (@mattnorris)
```

The name and email are optional, but recommended. The login is required.

If only the login is specified, do not put the login in parenthesis or prefix it with an @.

To use organizational approval, you must use the TOML file format to specify the members of each team. This format looks like:
```
[people]
    [people.bob]
        name = "Bob Bobson"
        email = "bob@lgtm.co"
        login = "bob"
    [people.fred]
        name = "Fred Fredson"
        email = "fred@lgtm.co"
        login = "fred"
    [people.jon]
        name = "Jon Jonson"
        email = "jon@lgtm.co"
        login = "jon"
    [people.ralph]
        name = "Ralph Ralphington"
        email = "ralph@lgtm.co"
        login = "ralph"
    [people.george]
        name = "George McGeorge"
        email = "george@lgtm.co"
        login = "george"

[org]
    [org.cap]
        people = [
            "bob",
            "fred",
            "jon"
        ]

    [org.iron]
        people = [
            "ralph",
            "george"
        ]
```

If no MAINTAINERS file is present, LGTM will fall back to using a github team with the name MAINTAINERS. You cannot enable organizational approval without a MAINTAINERS file.

#### Approval Management

In order for LGTM to pass its status check there must be a number of valid approvals. This is controlled by the `approvals` property in the .lgtm file:
```
approvals = 2
```

If this property is not specified, it defaults to the value 2.

The message that approvers use to signal that they have granted approval is controlled by the `pattern` property in the .lgtm file:
```
pattern = '(?i)LGTM'
```

If no pattern is specified, `(?i)LGTM` is used as the pattern.

##### Simple approval

By default, any approver on the list
```
approval_algorithm = "simple"
```

If you want to prevent the creator of a pull request from approving their own pull request, set the `self_approval_off` property in the .lgtm file:

```
self_approval_off = true
```

This property defaults to false; i.e., if you are on the approver list, you can approve a pull request that you create.

##### Organizational approval

Organizational approval means that at most one person from a specified organizational team can approve a pull request.

If self-approval is disabled (`self_approval_off = true`), then no one on the same team as the author can approve the pull request.

It's possible to configure organizations so that a person is on multiple teams, but it is not recommended, as the team assignment recognized by LGTM is undefined.

To enable organizational approval, add the following line to the .lgtm file:

```
approval_algorithm = "org"
```
