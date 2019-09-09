# Azure

## Azure specific configuration items

* DirectoryID
* ApplicationID
* ApplicationSecret
* SubscriptionID
* ResourceFilter - Use the resourceFilter setting to limit the resources from which metrics are collected. Otherwise, **ALL** metrics from **ALL** resources will be collected. The suggested method would be to add a tag to each resource from which to collect metrics and use a filter expression such as `tagname eq 'circonus' and tagValue eq 'enabled'`
* CloudName - default "AzurePublicCloud"
* UserAgent - default "circonus-cloud-agent"
* Interval - collection interval in minutes, default 5 [>=default]

Example configuration: `circonus-cloud-agentd --enable-azure --azure-example-conf=yaml`

## Setting up application in Azure

1. Login to the Azure portal
1. [Create application](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#create-an-azure-active-directory-application)
    1. name the application (circonus-cloud-agent)
    1. ~~assign _Monitoring Reader_ access role policy~~
    1. add reader role to subscription for application
1. collect application configuration information
    1. directory id
    1. application id
    1. application secret
    1. subscription id

## Bare minimum configuration

```yaml
---
id:
azure:
  directory_id:
  application_id:
  application_secret:
  subscription_id:
circonus:
  key:
```
