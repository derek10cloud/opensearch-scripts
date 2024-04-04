# opensearch search performance test

Script for opensearch (AWS Managed)

## environment variable

| Environment Variable | Example                               | Default Value  | Note |
| -------------------- | ------------------------------------- | -------------- | ---- |
| OPENSEARCH_ENDPOINT  | https://opensearch-domain.region.com/ |                |      |
| AWS_REGION           | ap-northeast-2                        | ap-northeast-2 |      |

## How to run

```shell
aws-vault exec {AWS Profile alias} -- go run main.go --index_name {index name} ...
```

You can check your {AWS Profile} from `~/.aws/config`
