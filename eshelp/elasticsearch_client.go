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
	client, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil {
		panic(err)
	}

	return &ElasticSearchClient{client, context.Background()}
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
