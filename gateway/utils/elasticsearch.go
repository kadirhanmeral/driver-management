package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"go-api-gateway/config"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	ESClient *elasticsearch.Client
	ESIndex  string
)

func InitElasticsearch(cfg *config.Config) {
	if cfg.ElasticsearchURL == "" {
		log.Println("ELASTICSEARCH_URL not set, skipping Elasticsearch initialization")
		return
	}

	ESIndex = cfg.ElasticsearchIndex
	if ESIndex == "" {
		ESIndex = "gateway-logs"
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchURL},
	}
	var err error
	ESClient, err = elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("Error creating the ES client: %s", err)
		return
	}

	res, err := ESClient.Info()
	if err != nil {
		log.Printf("Error connecting to Elasticsearch: %s", err)
		return
	}
	defer res.Body.Close()
	log.Println("Connected to Elasticsearch")
}

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	TraceID      string    `json:"traceId"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"statusCode"`
	ClientIP     string    `json:"clientIP"`
	UserAgent    string    `json:"userAgent"`
	RequestBody  string    `json:"requestBody,omitempty"`
	ResponseBody string    `json:"responseBody,omitempty"`
	LatencyMs    int64     `json:"latencyMs"`
}

func SendLogToES(entry LogEntry) {
	if ESClient == nil {
		return
	}
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %s", err)
		return
	}

	req := esapi.IndexRequest{
		Index:   ESIndex,
		Body:    bytes.NewReader(data),
		Refresh: "false",
	}

	go func() {
		res, err := req.Do(context.Background(), ESClient)
		if err != nil {
			log.Printf("Error indexing log: %s", err)
			return
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("Error indexing log: %s", res.String())
		}
	}()
}
