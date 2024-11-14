package weeder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"news-weeder/internal/weeder"
	"news-weeder/internal/weeder/redis"
)

func TestRedisSearch(t *testing.T) {
	knnConfig := redis.RedisConfig{
		Address: "localhost:6379",
		Index:   "test-vec-index",
		KNN:     5,
		DIM:     4,
	}

	knnSearch := redis.New(&knnConfig)

	doc1 := weeder.Document{
		ID:        "1",
		Title:     "ukraine war",
		Timestamp: time.Now(),
		Content:   "describes war escalation between russia and ukraine",
		Embedding: []float32{0.1, 0.2, 0.3, 0.4},
	}

	doc2 := weeder.Document{
		ID:        "2",
		Title:     "london library",
		Timestamp: time.Now(),
		Content:   "this article about the great london library",
		Embedding: []float32{0.5, 0.6, 0.7, 0.8},
	}

	t.Run("Search Semantic", func(t *testing.T) {
		err := knnSearch.Weeder.CreateSchema()
		assert.NoError(t, err)

		for _, doc := range []weeder.Document{doc1, doc2} {
			_ = knnSearch.Weeder.Append(&doc)
		}

		params1 := weeder.SearchParams{Limit: 1, DaysOffset: 2, Vector: []float32{0.5, 0.6, 0.7, 0.8}}
		result1, err := knnSearch.Weeder.Search(&params1)
		assert.NoError(t, err)
		assert.NotEmpty(t, result1)

		params2 := weeder.SearchParams{Limit: 1, DaysOffset: 2, Vector: []float32{0.1, 0.2, 0.3, 0.4}}
		result2, err := knnSearch.Weeder.Search(&params2)
		assert.NoError(t, err)
		assert.NotEmpty(t, result2)
	})
}
