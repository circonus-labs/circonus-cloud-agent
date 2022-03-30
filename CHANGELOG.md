# v0.3.2

* upd: enable using IAM role attached to an ec2 instance rather than explicit credential file or keys

# v0.3.1

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

# v0.3.0

* add: custom collector to walk through ec2 list of ebs volumes and collect from each one individually
* upd: dimensions `map[string]string`
* upd: dump _only_ on when `trace_metrics` setting true
* fix: lint warnings

# v0.2.1

* add: AWS/DX

# v0.2.0

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
* upd: add additional aws colleciton logging
* upd: adjust aws collection interval delta calculation
* upd: stricter linting
* doc: update install instructions (correct repository for releases)

# v0.1.1

* fix: release to circonus-labs repository

# v0.1.0

* Initial release
