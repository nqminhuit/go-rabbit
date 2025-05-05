package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"server/utils"
	"time"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

const INGEST_PIPELINE_NAME = "ingest_reports"

type OpenSearchClient struct {
	Client    *opensearchapi.Client
	IndexName string
}

func (o *OpenSearchClient) existIndex(indexName string) bool {
	req := opensearchapi.IndicesExistsReq{
		Indices: []string{indexName},
	}
	res, err := o.Client.Client.Do(context.Background(), req, nil)
	defer utils.Close(res.Body)
	utils.LogOnError(err, "Failed to check for existing index")

	if res.StatusCode == 200 {
		slog.Info("Index existed", "indexName", indexName)
		return true
	}
	return false
}

func (o *OpenSearchClient) createIndex(indexName string) {
	if o.existIndex(indexName) {
		return
	}
	req := opensearchapi.IndicesCreateReq{
		Index: indexName,
		Body: bytes.NewReader([]byte(`
		{
			"settings": {
				"index": {
					"codec": "zstd",
					"refresh_interval": "30s",
					"number_of_shards": 1,
					"number_of_replicas": 0
				}
			},
			"mappings": {
				"properties": {
					"data_id": {"type": "keyword", "ignore_above": 256},
					"batch_id": {"type": "keyword", "ignore_above": 256},
					"file_info": {
						"type": "object",
						"properties": {
							"display_name": {
								"type": "match_only_text",
								"fields": {"keyword": {"type": "keyword", "ignore_above": 256}}
							},
							"file_type": {"type": "keyword", "ignore_above": 256},
							"file_type_description": {"type": "keyword", "ignore_above": 256},
							"sha256": {"type": "keyword", "ignore_above": 256},
							"file_type_id": {"type": "keyword", "ignore_above": 256},
							"md5": {"type": "keyword", "ignore_above": 256},
							"sha1": {"type": "keyword", "ignore_above": 256},
							"sha512": {"type": "keyword", "ignore_above": 256},
							"type_category": {"type": "keyword", "ignore_above": 256}
						}
					},
					"process_info": {
						"type": "object",
						"properties": {
							"processing_time": {"type": "long"},
							"result": {"type": "keyword", "ignore_above": 256},
							"source": {"type": "keyword", "ignore_above": 256},
							"user_agent": {"type": "keyword", "ignore_above": 256},
							"username": {"type": "keyword", "ignore_above": 256},
							"post_processing": {"type": "flat_object"},
							"archive_handling_details": {"type": "flat_object"},
							"processing_time_details": {"type": "object"}
						}
					},
					"request_type": {"type": "keyword", "ignore_above": 256},
					"scan_results": {
						"type": "object",
						"properties": {
							"scan_all_result_a": {"type": "keyword", "ignore_above": 256},
							"scan_details": {"type": "flat_object"},
							"start_time": {"type": "date"},
							"total_avs": {"type": "long"},
							"total_time": {"type": "long"}
						}
					},
					"parent_path": {"type": "flat_object"},
					"insightsti_info": {"type": "flat_object"},
					"opswatfilescan_info": {"type": "flat_object"},
					"coo_info": {"type": "flat_object"},
					"sbom_info": {"type": "flat_object"},
					"vulnerability_info": {"type": "flat_object"},
					"yara_info": {"type": "flat_object"},
					"extraction_info": {"type": "flat_object"},
					"dlp_info": {"type": "flat_object"},
					"filetype_info": {"type": "flat_object"},
					"reputation_info": {"type": "flat_object"},
					"batch_files": {"type": "flat_object"},
					"extracted_files": {"type": "flat_object"}
				}
			}
		}`)),
	}
	res, err := o.Client.Client.Do(context.Background(), req, nil)
	defer utils.Close(res.Body)
	utils.LogOnError(err, "Failed to create index")
	slog.Info("Created", "index-info", res)
}

func (o *OpenSearchClient) createIngestPipeline() {
	req := opensearchapi.IngestCreateReq{
		PipelineID: INGEST_PIPELINE_NAME,
		Body: bytes.NewReader([]byte(`
		{
			"processors": [{
				"remove": {
					"field": [
						"filetype_info.file_info",
						"filetype_info.result_template_hash",
						"process_info.post_processing.copy_move_destination",
						"process_info.post_processing.converted_to",
						"process_info.post_processing.converted_destination",
						"process_info.post_processing.result_template_hash",
						"process_info.post_processing.sanitized_file_info",
						"process_info.post_processing.cdr_wait_time"
					],
					"ignore_failure": true
				}
			}]
		}`)),
	}
	res, err := o.Client.Client.Do(context.Background(), req, nil)
	defer utils.Close(res.Body)
	utils.LogOnError(err, "Failed to create ingest pipeline")
	slog.Info("Created ingest pipeline", "response", res)
}

func (o *OpenSearchClient) AddToBulk(indexer *opensearchutil.BulkIndexer, dataId string, data *bytes.Reader) {
	err := (*indexer).Add(context.Background(), opensearchutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: dataId,
		Body:       data,
		OnFailure: func(_ context.Context, _ opensearchutil.BulkIndexerItem, resp opensearchapi.BulkRespItem, err error) {
			if err != nil {
				slog.Error(err.Error())
			} else {
				slog.Error("Response failed from OpenSearch", "type", resp.Error.Type, "reason", resp.Error.Reason)
			}
		},
	})
	utils.LogOnError(err, "Failed to add item to indexer")

	// err = (*indexer).Close(context.Background())
	// utils.LogOnError(err, "Failed to close indexer")

	stats := (*indexer).Stats()
	if stats.NumFailed > 0 {
		slog.Error("There were failed documents", "total-flushed", stats.NumFlushed, "total-failed", stats.NumFailed)
	}
}

func ConnectToOpenSearch(coreIndexName string, username string, password string, addresses []string) *OpenSearchClient {
	slog.Info("Connecting to opensearch...")
	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport:           &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			Addresses:           addresses,
			Username:            username,
			Password:            password,
			RetryOnStatus:       []int{502, 503, 504, 429},
			RetryBackoff:        func(i int) time.Duration { return time.Duration(i) * 3 * time.Second },
			MaxRetries:          30,
			CompressRequestBody: true,
		},
	})
	utils.FailOnError(err, "Failed to connect to OpenSearch")
	info, err := client.Info(context.Background(), &opensearchapi.InfoReq{})
	utils.FailOnError(err, "Failed to check for OpenSearch info")
	slog.Info("Successfully connect to OpenSearch", "version", info.Version.Number, "name", info.Name)

	o := &OpenSearchClient{
		Client: client,
	}

	o.createIngestPipeline()
	o.createIndex(coreIndexName)
	return o
}
