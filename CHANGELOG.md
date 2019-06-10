# v0.0.1-alpha.5

* dev release: AWS MVP, Azure MVP, GCP MVP
* fix: check create - typo (asynch_metrics) causes aggregation on poll interval if not set to true
* fix: close connections after submit
* upd: svc log messages
* upd: pass interval to collectors
* upd: tick every min, run on time since last run >= interval
* upd: one common date range for timeseries requests
* upd: refactor submitting into collectors
* upd: refactor collector to include common date range for requests
* upd: refactor collector for submitting metrics per item
* upd: refactor collect method, submit metrics per instance
* fix: close all connections, api clients leaving Ks of https conns in ESTAB
* add: duration logging
* add: logging
* upd: clarify variable names
* upd: queue items from api calls, processing per item causes auth token timeouts
* fix: remove metric_kind and metric_type stream tags
* doc: update gcp setup steps
* upd: relax test configuration ignoring patterns

# v0.0.1-alpha.4

* dev release: AWS MVP, Azure MVP, GCP WIP
* gcp: refactor
* gcp: break on iter.Next() errors (503s and auth token failures after retrieving auth token...)
* doc: update comments
* upd: use logshim for api logging through zerolog
* upd: log check bundle that will be used
* upd: dependency (go-apiclient->retryablehttp for logging via zerolog)
* add: service interface and scan method stubs
* aws: rename 'creds' -> 'aws' in configuration
* doc: readme updates (main, azure, aws)
* upd: copyright year (all relevant files)
* upd: azure, refactor getMericData
* upd: base, add recconoiter metric type constants
* add: common error metric for exposing collection errors to users in their checks

# v0.0.1-alpha.3

* dev release
* azure: wip - initial
* aws: updates for base refactor
* base: refactor metric and tag handling, merge into check instance

# v0.0.1-alpha.2

* dev release
* aws: add es, elasticmapreduce, elasticbeanstalk, s3, rds, vpc (natgateway, transitgateway), lambda, kms
* aws: doc upd comments in collectors/metrics submitter method
* aws: upd copy base tags slice

# v0.0.1-alpha.1

* dev release
* aws: add ability to disable individual metrics
* aws: upd services requiring meta data, auto disable on access denied errors where extra permissions are required (retried every hour, so can be dynamically fixed by updating creds)
* aws: upd logging of aws specific errors for meta data requests
* aws: add elasticache
* aws: upd emit full config example for --aws-example-config

# v0.0.1-alpha

* dev release
* aws: applicationelb, cloudfront, dynamodb, ebs, ec2, ec2autoscaling, ec2spot, ecs, efs, elasticinterface, elastictranscoder, elb, networkelb, route53, route53resolver, sns, sqs

# v0.0.0

* skeleton development
