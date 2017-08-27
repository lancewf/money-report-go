package eshelp

import (
	"context"

	elastic "gopkg.in/olivere/elastic.v5"
)

type ElasticSearchClient struct {
	client *elastic.Client
	ctx    context.Context
}

func NewElasticSearchClient(esUrl string) *ElasticSearchClient {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetURL(esUrl))
	if err != nil {
		panic(err)
	}

	return &ElasticSearchClient{client, context.Background()}
}

func (client *ElasticSearchClient) IndexExists(indexName string) bool {
	exists, err := client.client.IndexExists().Index([]string{indexName}).Do(client.ctx)
	if err != nil {
		panic(err)
	}
	return exists
}

func (client *ElasticSearchClient) CreateMapping(indexName string, indexTypeName string, mappingText string) (bool, error) {
	putMappingResponse, error := client.client.PutMapping().
		Index(indexName).Type(indexTypeName).
		BodyString(mappingText).
		Do(client.ctx)
	if error != nil {
		panic(error)
	}
	return putMappingResponse.Acknowledged, error
}

func (client *ElasticSearchClient) CreateIndex(indexName string) {
	// Create an index
	client.client.CreateIndex(indexName).Do(client.ctx)
}

func (client *ElasticSearchClient) CreateDoc(indexName string, indexType string, body interface{}) {
	client.client.Index().
		Index(indexName).
		Type(indexType).
		BodyJson(body).
		Refresh("true").
		Do(client.ctx)
}
