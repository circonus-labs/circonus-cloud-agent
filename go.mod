module github.com/circonus-labs/circonus-cloud-agent

go 1.12

require (
	cloud.google.com/go/monitoring v1.16.3
	github.com/Azure/azure-sdk-for-go v33.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.1
	github.com/Azure/go-autorest/autorest/adal v0.6.0
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/alecthomas/units v0.0.0-20210927113745-59d0afb8317a
	github.com/aws/aws-sdk-go v1.23.17
	github.com/circonus-labs/go-apiclient v0.7.23
	github.com/golang/protobuf v1.5.3
	github.com/pelletier/go-toml v1.9.4
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	golang.org/x/oauth2 v0.15.0
	golang.org/x/sync v0.4.0
	golang.org/x/sys v0.15.0
	google.golang.org/api v0.149.0
	google.golang.org/genproto v0.0.0-20231016165738-49dd2c1f3d0b
	google.golang.org/genproto/googleapis/api v0.0.0-20231016165738-49dd2c1f3d0b
	gopkg.in/yaml.v2 v2.4.0
)

// replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
