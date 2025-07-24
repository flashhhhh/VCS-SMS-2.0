package elasticsearch

import (
	"bytes"
	"context"
	"errors"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/flashhhhh/pkg/env"
)

type ElasticsearchClient interface {
	Index(ctx context.Context, index string, body []byte) (error)
	Search(ctx context.Context, index string, buf bytes.Buffer) (*esapi.Response, error)
}

type elasticsearchClient struct {
	es *elasticsearch.Client
}

func NewElasticsearchClient(es *elasticsearch.Client) ElasticsearchClient {
	return &elasticsearchClient{
		es: es,
	}
}

func (esc *elasticsearchClient) Index(ctx context.Context, index string, data []byte) (error) {
	req := esapi.IndexRequest{
		Index:   env.GetEnv("ES_NAME", "ping_status"),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(ctx, esc.es)
	if err != nil {
		return errors.New("can't send request to ES")
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("Error response from ES: " + res.String())
	}

	return nil
}

func (esc *elasticsearchClient) Search(ctx context.Context, index string, buf bytes.Buffer) (*esapi.Response, error) {
	resp, err := esc.es.Search(
		esc.es.Search.WithContext(ctx),
		esc.es.Search.WithIndex(index),
		esc.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, errors.New("can't send search request to ES")
	}
	if resp.IsError() {
		defer resp.Body.Close()
		return nil, errors.New("Error response from ES: " + resp.String())
	}
	return resp, nil
}