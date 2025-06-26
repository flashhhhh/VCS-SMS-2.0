package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	eslib "server_administration_service/infrastructure/elasticsearch"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/flashhhhh/pkg/logging"
)

type ServerInfoRepository interface {
	GetServerMeanOnTimeRatio(startTime, endTime string) (float64, error)
}

type serverInfoRepository struct {
	es *elasticsearch.Client
}

func NewServerInfoRepository(es *elasticsearch.Client) ServerInfoRepository {
	return &serverInfoRepository{
		es: es,
	}
}

func (r *serverInfoRepository) GetServerMeanOnTimeRatio(startTime, endTime string) (float64, error) {
	startTimeInt64, err := time.Parse(time.RFC3339, startTime)
	if (err != nil) {
		logging.LogMessage("server_administration_service", "Start time is not valid", "ERROR")
		return 0, err
	}

	endTimeInt64, err := time.Parse(time.RFC3339, endTime)
	if (err != nil) {
		logging.LogMessage("server_administration_service", "End time is not valid", "ERROR")
		return 0, err
	}

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"Timestamp": map[string]interface{}{
					"gte": startTime,
					"lte": endTime,
				},
			},
		},
		"aggs": map[string]interface{}{
			"id_bucket": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "ID",
					"size": 10000,
				},
				"aggs": map[string]interface{}{
					"on_ping": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"Status.keyword": "On",
							},
						},
						"aggs": map[string]interface{}{
							"total_on_ping_time": map[string]interface{}{
								"sum": map[string]interface{}{
									"field": "Timestamp",
								},
							},
						},
					},
					"off_ping": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"Status.keyword": "Off",
							},
						},
						"aggs": map[string]interface{}{
							"total_off_ping_time": map[string]interface{}{
								"sum": map[string]interface{}{
									"field": "Timestamp",
								},
							},
						},
					},
					"last_ping": map[string]interface{}{
						"top_hits": map[string]interface{}{
							"size": 1,
							"sort": []map[string]interface{}{
								{
									"Timestamp": map[string]interface{}{
										"order": "desc",
									},
								},
							},
							"_source": map[string]interface{}{
								"includes": []string{"Status"},
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logging.LogMessage("server_administration_service", "Failed to encode query to buffer. Err: " + err.Error(), "ERROR")
		return 0, err
	}

	resp, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("ping_status"),
		r.es.Search.WithBody(&buf),
	)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get Elasticsearch response. Err: " + err.Error(), "ERROR")
		return 0, err
	}

	defer resp.Body.Close()

	var answer map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		logging.LogMessage("server_administration_service", "Failed to decode Elasticsearch query's result, err: " + err.Error(), "ERROR")
		return 0, err
	}

	buckets := answer["aggregations"].(map[string]interface{})["id_bucket"].(map[string]interface{})["buckets"].([]interface{})
	meanOnTimeRatio := 0.0
	
	for _, s := range buckets {
		server := s.(map[string]interface{})
		
		last_ping := server["last_ping"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})
		last_status := last_ping["_source"].(map[string]interface{})["Status"].(string)
		
		doc_count := int(server["doc_count"].(float64))
		first_status := last_status

		if (doc_count % 2 == 0) {
			switch last_status {
			case "On":
				first_status = "Off"
			case "Off":
				first_status = "On"
			}
		}

		off_ping := server["off_ping"].(map[string]interface{})
		total_off_ping_time := int64(off_ping["total_off_ping_time"].(map[string]interface{})["value"].(float64))
		
		on_ping := server["on_ping"].(map[string]interface{})
		total_on_ping_time := int64(on_ping["total_on_ping_time"].(map[string]interface{})["value"].(float64))
		
		if last_status == "On" {
			total_off_ping_time += endTimeInt64.UnixNano() / int64(time.Millisecond)
		}

		if first_status == "Off" {
			total_on_ping_time += startTimeInt64.UnixNano() / int64(time.Millisecond)
		}

		total_time := (endTimeInt64.UnixNano() - startTimeInt64.UnixNano()) / int64(time.Millisecond)

		on_time_ratio := float64(total_off_ping_time-total_on_ping_time) / float64(total_time) * 100
		meanOnTimeRatio += on_time_ratio
	}

	meanOnTimeRatio /= float64(len(buckets))

	return float64(meanOnTimeRatio), nil
}

func main() {
	esAddress := "http://localhost:9200"
	es := eslib.ConnectES(esAddress)

	serverInfoRepository := NewServerInfoRepository(es)
	meanOnTimeRatio, _ := serverInfoRepository.GetServerMeanOnTimeRatio("2025-06-24T00:00:00Z", "2025-06-24T23:59:59Z")
	fmt.Printf("%.6f\n", meanOnTimeRatio)
}