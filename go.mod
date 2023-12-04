module github.com/circonus-labs/circonus-cloud-agent

go 1.12

require (
	cloud.google.com/go/monitoring v1.0.0
	github.com/Azure/azure-sdk-for-go v33.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.1
	github.com/Azure/go-autorest/autorest/adal v0.6.0
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
	github.com/aws/aws-sdk-go v1.23.17
	github.com/circonus-labs/go-apiclient v0.7.15
	github.com/golang/protobuf v1.5.2
	github.com/pelletier/go-toml v1.9.4
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.15.0
	google.golang.org/api v0.57.0
	google.golang.org/genproto v0.0.0-20210921142501-181ce0d877f6
	gopkg.in/yaml.v2 v2.4.0
)

// replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
