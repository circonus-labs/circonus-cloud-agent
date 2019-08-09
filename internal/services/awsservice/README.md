# AWS CloudWatch Metrics

## Status

Initial draft implementation -- **alpha** -- a MVP at this point.

## Supported

* ApplicationELB
* CloudFront
* DynamoDB
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

## Configuration

### AWS

1. Create an IAM group e.g. `circonus-cloud-agent` (if one does not already exist) with the following permissions:
   1. Required: `CloudWatchReadOnlyAccess` - to retrieve metrics
   1. Optional: `AmazonEC2ReadOnlyAccess` - to get list of instances, if EC2 metrics are desired
   1. Optional: `AmazonElastiCacheReadOnlyAccess` - to get list of cache clusters, if ElastiCache metrics are desired
1. Create an IAM user e.g. `circonus-cloud-agent` (ensure it is a member of the aforementioned group), make a note of the Access Key Credentials (access key id and secret access key or role and optionally credentials file for local credentials)

> NOTE: the agent can run in a shared mode with multiple sets of credentials or a local mode where the credentials are kept in an external location and _roles_ are used to identify which credentials to use.

### Circonus

1. Use Circonus UI to create or identify an API Token to use
1. Create a configuration file (`etc/aws.d/...` with some unique name), place the aws and circonus credentials where appropriate and configure what metrics to retrieve
1. Start/restart `circonus-cloud-agentd`, ensure `--enable-aws` is on command line, environment variable is set, or aws is enabled in main configuration file

### Example -- for shared

[example-config.yaml](example-config.yaml) is a **full** example. Note that a simple basic configuration for a single service such as EC2 (provided the default metrics are acceptable) would be:

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

### Example -- for local

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
