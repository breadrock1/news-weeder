package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/rueidis"
	"news-weeder/internal/weeder"
)

type KNNSearch struct {
	config *RedisConfig
	rsConn *redis.Client
}

func New(config *RedisConfig) *weeder.Weeder {
	redisOpts := &redis.Options{Addr: config.Address}
	conn := redis.NewClient(redisOpts)
	weederClient := &KNNSearch{
		config: config,
		rsConn: conn,
	}

	return &weeder.Weeder{Weeder: weederClient}
}

func (rs *KNNSearch) CreateSchema() error {
	dim := rs.config.DIM
	index := rs.config.Index
	createIndex := []interface{}{
		"FT.CREATE", index, "ON", "HASH",
		"PREFIX", "1", "doc:",
		"SCHEMA",
		"title", "TEXT",
		"timestamp", "NUMERIC", "SORTABLE",
		"embedding", "VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", dim, "DISTANCE_METRIC", "COSINE",
	}

	ctx := context.Background()
	err := rs.rsConn.Do(ctx, createIndex...).Err()
	if err != nil && err.Error() != "Index already exists" {
		return err
	}

	return nil
}

func (rs *KNNSearch) Delete(docId string) error {
	index := rs.config.Index
	key := fmt.Sprintf("doc:%s", docId)

	err := rs.rsConn.HDel(context.Background(), index, key).Err()
	return err
}

func (rs *KNNSearch) Append(doc *weeder.Document) error {
	key := fmt.Sprintf("doc:%s", doc.ID)
	hashData := rueidis.VectorString64(doc.Embedding)

	setQuery := []interface{}{
		"title", doc.Title,
		"timestamp", doc.Timestamp.Unix(),
		"embedding", hashData,
	}

	ctx := context.Background()
	err := rs.rsConn.HSet(ctx, key, setQuery...).Err()
	return err
}

func (rs *KNNSearch) Search(params *weeder.SearchParams) ([]*weeder.Document, error) {
	var filter string
	if params.DaysOffset < 1 {
		filter = "*"
	} else {
		offset := time.Now().Add(time.Duration(params.DaysOffset) * time.Hour * -1)
		filter = fmt.Sprintf("@timestamp:[%d inf+] ", offset.Unix())
	}

	knn := fmt.Sprintf("%s=>[KNN %d @embedding $vec AS score]", filter, rs.config.KNN)

	index := rs.config.Index
	embed := rueidis.VectorString64(params.Vector)

	searchQuery := []interface{}{
		"FT.SEARCH", index,
		knn,
		"PARAMS", "2", "vec", embed,
		"RETURN", "3", "score", "title", "timestamp",
		"LIMIT", "0", "10",
		"SORTBY", "score", "ASC",
		"DIALECT", "2",
	}

	ctx := context.Background()
	res, err := rs.rsConn.Do(ctx, searchQuery...).Result()
	if err != nil {
		return nil, err
	}

	docs, err := rs.extractSearchResult(res)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (rs *KNNSearch) extractSearchResult(foundResult interface{}) ([]*weeder.Document, error) {
	resMap := foundResult.(map[interface{}]interface{})
	total := resMap["total_results"].(int64)
	results := resMap["results"].([]interface{})

	foundedDocs := make([]*weeder.Document, total)
	for ind, resItem := range results {
		item := resItem.(map[interface{}]interface{})
		id := item["id"].(string)

		attrs := item["extra_attributes"].(map[interface{}]interface{})
		title := attrs["title"].(string)

		score := attrs["score"].(string)
		floatObj, err := strconv.ParseFloat(score, 32)
		if err != nil {
			return nil, err
		}

		var articleTimestamp int64
		timestamp := attrs["timestamp"].(string)
		articleTimestamp, err = strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			articleTimestamp = time.Now().Unix()
		}
		timeVal := time.Unix(articleTimestamp, 0)

		foundDoc := &weeder.Document{
			ID:        id,
			Title:     title,
			Timestamp: timeVal,
			Score:     float32(floatObj),
		}

		foundedDocs[ind] = foundDoc
	}

	return foundedDocs, nil
}

func (rs *KNNSearch) Connect() error {
	return nil
}

func (rs *KNNSearch) Close() error {
	return nil
}
