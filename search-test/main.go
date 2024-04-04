package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

var (
	endpoint = os.Getenv("OPENSEARCH_ENDPOINT") // e.g. https://opensearch-domain.region.com or Amazon OpenSearch Serverless endpoint
)

func main() {
	indexName := flag.String("index_name", "", "Index Name")
	aws_region := flag.String("region", "ap-northeast-2", "AWS Region")
	flag.Parse()

	ctx := context.Background()

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(*aws_region),
	)
	if err != nil {
		log.Fatal(err) // Do not log.fatal in a production ready app.
	}

	// Create an AWS request Signer and load AWS configuration using default config folder or env vars.
	signer, err := requestsigner.NewSignerWithService(awsCfg, "es")
	if err != nil {
		log.Fatal(err) // Do not log.fatal in a production ready app.
	}

	// Create an opensearch client and use the request-signer
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{endpoint},
		Signer:    signer,
	})
	if err != nil {
		log.Fatal("client creation err", err)
	}

	content := strings.NewReader(`{
		"size": 1
	}`)

	search := opensearchapi.SearchRequest{
		Index: []string{*indexName},
		Body:  content,
	}

	searchResponse, err := search.Do(context.Background(), client)
	fmt.Println("Response", searchResponse.String())
}
