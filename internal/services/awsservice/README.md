# AWS CloudWatch Metrics

## Supported AWS services

* ApplicationELB
* CloudFront
* DynamoDB
* DX
* EBS
* EC2
* EC2AutoScaling
* EC2Spot
* ECS
* EFS
* ElastiCache
* ElasticBeanstalk
* ElasticInterface
* ElasticMapReduce
* ElasticTranscoder
* ELB
* ES
* KMS
* Lambda
* NATGateway
* NetworkELB
* RDS
* Route53
* Route53Resolver
* S3
* SNS
* SQS
* TransitGateway

## Installation

1. Create a directory for the install. Suggested: `mkdir -p /opt/circonus/cloud-agent`
1. Download the [latest release](https://github.com/circonus-labs/circonus-cloud-agent/releases/latest)
1. Unpack the release in the directory created in first step
1. In this directory, create a config folder. Suggested: `mkdir /opt/circonus/cloud-agent/etc/aws.d`
1. Auto-create a service specific configuration template in the desired format (yaml, toml, or json).  Suggested: `sbin/circonus-cloud-agentd --enable-aws --aws-example-conf=yaml > /opt/circonus/cloud-agent/etc/aws.d/aws-config.yaml` 
    * Note, the `id` in the template is defaulted to the filename.  This should be changed to a name that will be unique across all cloud-agents used in Circonus
    * Add the [Circonus api key](#circonus)
    * Add the [AWS credentials](#aws-settings)
    * Update settings for the desired AWS services to be monitored
1. Setup as a system service or run in foreground ensuring that `--enable-aws` is specified

## Options

```sh
$  sbin/circonus-cloud-agentd -h
The Circonus Cloud Agent collects metrics from cloud infrastructures and fowards them to Circonus.

Usage:
  circonus-cloud-agent [flags]

Flags:
      --aws-conf-dir string         AWS configuration directory (default "/opt/circonus/cloud-agent/etc/aws.d")
      --aws-example-conf string     Show AWS config (json|toml|yaml) and exit
  -c, --config string               config file (default: circonus-cloud-agent.yaml|.json|.toml)
  -d, --debug                       [ENV: CCA_DEBUG] Enable debug messages
      --enable-aws                  Enable AWS metric collection client
  -h, --help                        help for circonus-cloud-agent
      --log-level string            [ENV: CCA_LOG_LEVEL] Log level [(panic|fatal|error|warn|info|debug|disabled)] (default "info")
      --log-pretty                  [ENV: CCA_LOG_PRETTY] Output formatted/colored log lines [ignored on windows]
  -V, --version                     Show version and exit
```

## Configuration

### AWS

> NOTE: the agent can run in a shared mode with multiple sets of credentials or a local mode where the credentials are kept in an external location and _roles_ are used to identify which credentials to use.

1. Create an IAM group e.g. `circonus-cloud-agent` (if one does not already exist) with the following permissions:
   1. Required: `CloudWatchReadOnlyAccess` - to retrieve metrics
   1. Optional: `AmazonEC2ReadOnlyAccess` - to get list of instances, if EC2 metrics are desired
   1. Optional: `AmazonElastiCacheReadOnlyAccess` - to get list of cache clusters, if ElastiCache metrics are desired
1. Create an IAM user e.g. `circonus-cloud-agent` (ensure it is a member of the aforementioned group), make a note of the Access Key Credentials (access key id and secret access key or role and save credentials file locally where circonus-cloud-agent will run, ensure it is accessible by whatever user the circonus-cloud-agent is configured to run as.)

### Circonus

1. Use Circonus UI to create or identify an API Token to use
1. Add the `key` to the config file under the `circonus` section

### AWS settings

For per-configuration file credentials (shared):

* `access_key_id`
* `secret_access_key`

or, for credentials in a local file:

* `role`
* `credentials_file`

### Example configuration

Minimum configuration (for EC2 service):

Credentials in configuration file:

```yaml
---

id: ...
aws:
  access_key_id: ...
  secret_access_key: ...
circonus:
  key: ...
period: basic
regions:
    - name: us-east-1
      services:
          - namespace: aws/EC2
            disabled: false
```

Shared credentials using roles:

```yaml
---

id: ...
aws:
  role: ...
  credentials_file: ...
circonus:
  key: ...
period: basic
regions:
    - name: us-east-1
      services:
          - namespace: aws/EC2
            disabled: false
```
