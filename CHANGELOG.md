# unreleased

## v0.3.3

* build(deps): bump google.golang.org/api from 0.149.0 to 0.152.0
* build(deps): bump github.com/Azure/go-autorest/autorest from 0.9.1 to 0.11.29
* build(deps): bump github.com/spf13/cobra from 1.2.1 to 1.8.0
* build(deps): bump github.com/rs/zerolog from 1.25.0 to 1.31.0
* build: add after hook for `grype` on generated sboms
* build: add .sbom for archive artifacts and update changelog config
* fix(lint): deprecated ioutil calls
* build(deps): update azure autorest/adal to v0.9.23 to address GO-2020-0017
* chore: update to go1.21
* fix: remove deprecated syntax
* build: add before hooks for `go mod tidy`, `govulncheck` and `golangci-lint`
* fix(lint): unused args and replacement monitoringpb module for gcp
* build(deps): bump golang.org/x/sync from 0.0.0-20210220032951-036812b2e83c to 0.5.0
* build(deps): bump cloud.google.com/go/monitoring from 1.0.0 to 1.16.3
* build(deps): bump github.com/circonus-labs/go-apiclient from 0.7.15 to 0.7.23
* build(deps): bump golang.org/x/oauth2 from 0.0.0-20210819190943-2bc19b11175f to 0.15.0
* fix(lint): various lint issues, typos, and struct member alignments

## v0.3.2

* upd: enable using IAM role attached to an ec2 instance rather than explicit credential file or keys

## v0.3.1

* fix: lint issues
* add: lint config
* add: lint workflow
* upd: remove log-level info msg
* bring select dependencies up-to-date
    * cloud.google.com/go/monitoring v1.0.0
    * github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
    * github.com/circonus-labs/go-apiclient v0.7.15
    * github.com/golang/protobuf v1.5.2
    * github.com/pelletier/go-toml v1.9.4
    * github.com/pkg/errors v0.9.1
    * github.com/rs/zerolog v1.25.0
    * github.com/spf13/cobra v1.2.1
    * github.com/spf13/viper v1.9.0
    * golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
    * golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
    * golang.org/x/sys v0.0.0-20210908233432-aa78b53d3365
    * google.golang.org/api v0.57.0
    * google.golang.org/genproto v0.0.0-20210921142501-181ce0d877f6
    * gopkg.in/yaml.v2 v2.4.0

## v0.3.0

* add: custom collector to walk through ec2 list of ebs volumes and collect from each one individually
* upd: dimensions `map[string]string`
* upd: dump _only_ on when `trace_metrics` setting true
* fix: lint warnings

## v0.2.1

* add: AWS/DX

## v0.2.0

* doc: update azure readme
* upd: clarify time delta calculation
* add: more aggressive ticker to start collections
* add: log collection duration
* upd: remove dead code
* fix: dynamic dimensions from call param
* fix: catch zero instances for a given region
* add: cloud vendor specific check submodule
* upd: add region debug message
* upd: reduce ticker drift to 2sec
* upd: promote duration message to Info
* upd: only print submission in debug mode
* upd: print response body as raw json
* fix: verify returned check bundle has a submission url
* upd: dependencies
* add: use service specific check submodules
* upd: add interval logging on client start
* upd: add additional aws collection logging
* upd: adjust aws collection interval delta calculation
* upd: stricter linting
* doc: update install instructions (correct repository for releases)

## v0.1.1

* fix: release to circonus-labs repository

## v0.1.0

* Initial release
