package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	"github.com/rs/zerolog/log"
)

var (
	indexName              = flag.String("index_name", "", "Index Name")
	aws_region             = flag.String("region", "ap-northeast-2", "AWS Region")
	request                = flag.Int("request", 3, "Request Number for each index query to get max,min,avg time")
	target_avg_search_time = flag.Int("target_avg_search_time", 1, "Target Average Search Time(ms)")
	url                    = flag.String("url", "", "endpoint of the OpenSearch cluster")
)

type Response struct {
	Took      int   `json:"took"`
	Timed_out bool  `json:"timed_out"`
	Shards    Shard `json:"_shards"`
}

type Shard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

func main() {
	// Parse the command-line flags
	flag.Parse()

	var (
		min int
		max int
		avg int
	)

	ctx := context.Background()
	client := makeOpensearchClient(ctx)
	for avg < *target_avg_search_time {
		min, max, avg = executeSearchQuery(client, *indexName, *request)
		// Print the header
		fmt.Println()
		fmt.Println()
		fmt.Println("===========================================================================Result===========================================================================")
		fmt.Println("============================================================================================================================================================")
		fmt.Println()
		fmt.Printf("Request Number: %d\n", *request)
		fmt.Printf("[Index: %s] Average Time: %d ms\n", *indexName, avg)
		fmt.Printf("[Index: %s] Minimum Search Time: %d ms\n", *indexName, min)
		fmt.Printf("[Index: %s] Maximum Time: %d ms\n", *indexName, max)
		fmt.Println()
		fmt.Println("============================================================================================================================================================")
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("===========================================================================Finish===========================================================================")
}

// makeOpensearchClient creates an opensearch client with the given context
func makeOpensearchClient(ctx context.Context) *opensearch.Client {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(*aws_region),
	)
	if err != nil {
		log.Error().Err(err).Msg("aws config error")
	}

	// Create an AWS request Signer and load AWS configuration using default config folder or env vars.
	signer, err := requestsigner.NewSignerWithService(awsCfg, "es")
	if err != nil {
		log.Error().Err(err).Msg("signer error")
	}

	// Create an opensearch client and use the request-signer
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{*url},
		Signer:    signer,
	})
	if err != nil {
		log.Error().Err(err).Msg("client creation err")
	}
	return client
}

// executeSearchQuery executes search queries on the given index
func executeSearchQuery(client *opensearch.Client, indexName string, request int) (min, max, avg int) {
	tookSlice := make([]int, request)
	for i := 0; i < request; i++ {
		var (
			resp *Response
		)
		content := strings.NewReader(`{
			"query": {
				"query_string": {
					"query": "*"
				}
			}
		}`)

		search := opensearchapi.SearchRequest{
			Index: []string{indexName},
			Body:  content,
		}

		searchResponse, err := search.Do(context.Background(), client)

		body, err := ioutil.ReadAll(searchResponse.Body)
		if body != nil {
			if err = json.Unmarshal(body, &resp); err != nil {
				log.Error().Err(err).Msg("unmashar error")
			}
		}

		// log.Info().Interface("body", searchResponse.Body).Msg("Body")
		// log.Info().Interface("response", resp).Msg("Response")
		// fmt.Println()
		// fmt.Printf("Took: %d\n", resp.Took)
		// fmt.Printf("Timed Out: %t\n", resp.Timed_out)
		// fmt.Printf("Total Shards: %d\n", resp.Shards.Total)
		// fmt.Printf("Successful Shards: %d\n", resp.Shards.Successful)
		// fmt.Printf("Skipped Shards: %d\n", resp.Shards.Skipped)
		// fmt.Printf("Failed Shards: %d\n", resp.Shards.Failed)
		tookSlice[i] = resp.Took
		time.Sleep(1 * time.Second)
	}
	min, max, avg = findMinMaxAverage(tookSlice)

	return min, max, avg
}

// findMinMaxAverage finds the minimum, maximum, and average values of the given slice
func findMinMaxAverage(slice []int) (min, max, avg int) {
	if len(slice) == 0 {
		return 0, 0, 0
	}

	min, max, total := slice[0], slice[0], 0

	for _, value := range slice {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
		total += value
	}

	avg = total / len(slice)

	return min, max, avg
}
