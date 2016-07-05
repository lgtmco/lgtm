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
Each repository managed by LGTM must have a .lgtm file in the root of the repository. This file provides the configuration that LGTM uses for the repository.

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
pattern = '(?i)LGTM\s*(\S*)'
```

If no pattern is specified, `(?i)LGTM\s*(\S*)` is used as the pattern.

If you want to use versioning, the capture group must be specified as part of the regular expression. If you don't, it can be left off.
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

#### Merging
LGTM has the ability to automatically merge a pull request when the branch is mergeable and all status checks (including optional ones) have passed. This feature is enabled on a repo-by-repo basis.

To enable merging, add the following line to the .lgtm file:

 ```
 do_merge = true
 ```

#### Versioning
LGTM has the ability to automatically place a version tag on a merge. This feature is enabled and configured on a repo-by-repo basis. If merging is not enabled, enabling versioning will do nothing.

To enable versioning, add the following line to the .lgtm file:
```
do_version = true
```

There are two kinds of versioning algorithms supported: semantic versioning and timestamps.

##### Semantic Versioning
The default behavior is to use Semantic Versioning. You can specify it explicitly with the following line in the .lgtm file:

 ```
 version_algorithm = "semver"
 ```

The semantic versioning standard is described at <http://semver.org/>.

Approval comments are used to provide the semantic version for the current pull request. By default, the version is specified after the LGTM string:

```
LGTM 1.0.0
```

will tag the merge commit with the version 1.0.0.

If multiple approvers for a pull commit specify a version, the highest version specified will be used.

If no approver specifies a version, or if the maximum version specified by an approver is less than any previously specified version tag, the version will be set to the highest specified version tag with the patch version incremented.


##### Timestamp Versioning

Timestamp versioning is enabled by adding the following line to the .lgtm file:

```
version_algorithm = "timestamp"
```

There are currently two options for timestamps, one based on the RFC 3339 standard and one specifying milliseconds since the epoch (Jan 1, 1970, 12:00:00AM UTC).

The default is the RFC 3339 format, but it can be explicitly specified by adding the following line to the .lgtm file:
```
version_format = "rfc3339"
```
The version tag will look like this:
```
2016-05-16T19.06.26Z
```
(colons aren't legal in git tags, so the colons in the RFC format have been replaced by periods.)

To enable millisecond versioning, add the following line to the .lgtm file:
```
version_format = "mills"
```
The version tag will look like this:
```
1463425671
```