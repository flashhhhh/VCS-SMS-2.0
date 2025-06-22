package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/flashhhhh/pkg/logging"
)

func ConnectES(dsn string) (*elasticsearch.Client) {
	client, connectES_err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{dsn},
	})

	if connectES_err != nil {
		logging.LogMessage("server_administration_service", "Error creating Elasticsearch client: " + connectES_err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("server_administration_service", "Connected to Elasticsearch at "+dsn, "INFO")
	return client
}

func CreateDocument(es *elasticsearch.Client, indexName string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return errors.New("Can't convert document to JSON")
	}

	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return errors.New("Can't send request to ES")
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("Error response from ES: " + res.String())
	}

	return nil
}

// BulkCreateDocuments adds multiple documents to Elasticsearch in a single bulk request.
// func BulkCreateDocuments(es *elasticsearch.Client, indexName string, docs []interface{}) error {
// 	var buf bytes.Buffer
// 	for _, doc := range docs {
// 		meta := map[string]interface{}{
// 			"index": map[string]interface{}{
// 				"_index": indexName,
// 			},
// 		}
// 		metaLine, err := json.Marshal(meta)
// 		if err != nil {
// 			return errors.New("Can't marshal bulk meta")
// 		}
// 		docLine, err := json.Marshal(doc)
// 		if err != nil {
// 			return errors.New("Can't marshal document")
// 		}
// 		buf.Write(metaLine)
// 		buf.WriteByte('\n')
// 		buf.Write(docLine)
// 		buf.WriteByte('\n')
// 	}

// 	req := esapi.BulkRequest{
// 		Body:    bytes.NewReader(buf.Bytes()),
// 		Refresh: "true",
// 	}

// 	res, err := req.Do(context.Background(), es)
// 	if err != nil {
// 		return errors.New("Can't send bulk request to ES")
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return errors.New("Bulk error response from ES: " + res.String())
// 	}

// 	return nil
// }