# Circonus Cloud Agent

Current status is a **work in progress** (WIP aka active development) of a minimum viable product (MVP) which will be refined based on feedback from initial/select users.

## Installation

1. Create a directory for the install: `mkdir -p /opt/circonus/cloud-agent`
1. Download the [latest release](https://github.com/maier/circonus-cloud-agent/releases/latest) (note: this is currently a private repository)
1. Unpack the release in the directory created in first step
1. Use link(s) below for configuration information for specific service(s)
1. Setup as a system service or run in foreground (e.g. for AWS after creating a configuration in `etc/aws.d` - `sbin/circonus-cloud-agentd --enable-aws`)

## Options

```sh
$  sbin/circonus-cloud-agentd -V
circonus-cloud-agent v0.0.1-alpha.3 - commit: cb3e4fa8773daf80557f2448ed85ffe88841349e, date: 2019-03-01T21:34:54Z, tag: v0.0.1-alpha.3

$  sbin/circonus-cloud-agentd -h
The Circonus Cloud Agent collects metrics from cloud infrastructures and fowards them to Circonus.

Usage:
  circonus-cloud-agent [flags]

Flags:
      --aws-conf-dir string         AWS configuration directory (default "/opt/circonus/cloud-agent/etc/aws.d")
      --aws-example-conf string     Show AWS config (json|toml|yaml) and exit
      --azure-conf-dir string       Azure configuration directory (default "/opt/circonus/cloud-agent/etc/azure.d")
      --azure-example-conf string   Show Azure config (json|toml|yaml) and exit
      --gcp-conf-dir string         GCP configuration directory (default "/opt/circonus/cloud-agent/etc/gcp.d")
      --gcp-example-conf string     Show GCP config (json|toml|yaml) and exit
  -c, --config string               config file (default: circonus-cloud-agent.yaml|.json|.toml)
  -d, --debug                       [ENV: CCA_DEBUG] Enable debug messages
      --enable-aws                  Enable AWS metric collection client
      --enable-azure                Enable Azure metric collection client
      --enable-gcp                  Enable GCP metric collection client
  -h, --help                        help for circonus-cloud-agent
      --log-level string            [ENV: CCA_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty                  [ENV: CCA_LOG_PRETTY] Output formatted/colored log lines [ignored on windows]
      --pipe-submits                [ENV: CCA_PIPE_SUBMITS] Pipe metric submissions to Circonus (experimental)
      --show-config string          Show config (json|toml|yaml) and exit
  -V, --version                     Show version and exit
```

For service specifc configurations see [links below](#services).

## Services

* [AWS](internal/services/awsservice/) - MVP
* [Azure](internal/services/azureservice) - MVP
* [GCP](internal/services/gcpservice) - MVP
