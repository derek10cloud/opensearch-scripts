# opensearch search performance test

Script for opensearch (AWS Managed)

## environment variable

| Arguments              | Example                               | Default Value  | optional | Note                                                                      |
| ---------------------- | ------------------------------------- | -------------- | -------- | ------------------------------------------------------------------------- |
| url                    | https://opensearch-domain.region.com/ |                | False    | Endpoint of the OpenSearch cluster                                        |
| index_name             | 2024.04.01-log                        |                | False    | Index Name of cluster to perform search test                              |
| region                 | ap-northeast-2                        | ap-northeast-2 | True     | AWS region of cluster                                                     |
| request                | 10                                    | 3              | True     | Request(every second) Number for each index query to get max,min,avg time |
| target_avg_search_time | 100                                   | 1              | True     | Target Average Search Time(ms)                                            |

## How to run

```shell
usage: aws-vault exec {AWS Profile alias} -- go run main.go --url URL [--region REGION] [--request REQUEST]
                                                            [--target_avg_search_time TARGET_AVG_SEARCH_TIME]
                                                            --index_name INDEX_NAME

necessary arguments:
  --index_name INDEX_NAME      Index Name of cluster to perform search test
  --url        URL             Endpoint of the OpenSearch cluster

optional arguments:
  --region                   REGION                     AWS region of cluster (default ap-northeast-2)
  --request                  REQUEST                    Request Number for each index query to get max,min,avg time, this requst will be performed in every second (default 3)
  --target_avg_search_time   TARGET_AVG_SEARCH_TIME     Target Average Search Time(ms) (default 1)
# e.g)
```

You can check your {AWS Profile} from `~/.aws/config`
