package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kusnadi8605/news/config"
	"github.com/kusnadi8605/news/entity"
	"github.com/kusnadi8605/news/repository"
	"github.com/sirupsen/logrus"
)

type NewsUsecaseImpl struct {
	Repo repository.NewsRepository
}

func NewNewsUsecase(repo repository.NewsRepository) NewsUsecase {
	return &NewsUsecaseImpl{Repo: repo}
}

// getFromCache mencoba mengambil data dari Redis
func getFromCache[T any](ctx context.Context, key string, target *T) (bool, error) {
	data, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err // Cache miss
	}

	if err := json.Unmarshal([]byte(data), target); err != nil {
		logrus.WithError(err).Warn("Failed to unmarshal cached data")
		return false, err
	}

	return true, nil
}

// saveToCache menyimpan data ke Redis dengan TTL
func saveToCache(ctx context.Context, key string, value any, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal data for Redis")
		return
	}

	if err := config.RedisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		logrus.WithError(err).Error("Failed to save data to Redis")
	}
}

// invalidateCache menghapus cache di Redis
func invalidateCache(ctx context.Context, key string) {
	if err := config.RedisClient.Del(ctx, key).Err(); err != nil {
		logrus.WithError(err).Error("Failed to delete cache from Redis")
	}
}

func (u *NewsUsecaseImpl) GetAllNews(ctx context.Context) (*[]entity.News, error) {
	cacheKey := "news:all"

	// Cek cache di Redis
	var newsList *[]entity.News
	if found, _ := getFromCache(ctx, cacheKey, &newsList); found {
		logrus.Info("Cache hit: Fetching news from Redis")
		return newsList, nil
	}

	logrus.Info("Cache miss: Fetching news from database")

	// Fetch dari database
	newsList, err := u.Repo.GetAllNews()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch news from database")
		return nil, err
	}

	// Simpan ke cache
	saveToCache(ctx, cacheKey, newsList, 10*time.Second)

	logrus.Info("News stored in Redis for caching")
	return newsList, nil
}

func (u *NewsUsecaseImpl) GetNewsByID(ctx context.Context, id int) (*entity.News, error) {
	cacheKey := fmt.Sprintf("news:%d", id)

	var news entity.News
	if found, _ := getFromCache(ctx, cacheKey, &news); found {
		logrus.WithField("news_id", id).Info("Cache hit: Fetching news from Redis")
		return &news, nil
	}

	logrus.WithField("news_id", id).Info("Cache miss: Fetching news from database")

	// Fetch dari database
	newsPtr, err := u.Repo.GetNewsByID(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch news from database")
		return nil, err
	}

	// Simpan ke cache jika data ditemukan
	if newsPtr != nil {
		saveToCache(ctx, cacheKey, newsPtr, 10*time.Second)
		logrus.WithField("news_id", id).Info("News stored in Redis for caching")
	}

	return newsPtr, nil
}

func (u *NewsUsecaseImpl) CreateNews(ctx context.Context, news *entity.News) (int64, error) {
	return u.Repo.CreateNews(news)
}

func (u *NewsUsecaseImpl) UpdateNews(ctx context.Context, news *entity.News) error {
	cacheKey := fmt.Sprintf("news:%d", news.ID)

	if err := u.Repo.UpdateNews(news); err != nil {
		logrus.WithError(err).Error("Failed to update news in database")
		return err
	}

	// Hapus cache terkait agar data terbaru digunakan
	invalidateCache(ctx, cacheKey)
	logrus.WithField("news_id", news.ID).Info("Cache invalidated after news update")
	return nil
}

func (u *NewsUsecaseImpl) DeleteNews(ctx context.Context, id int) error {
	cacheKey := fmt.Sprintf("news:%d", id)

	if err := u.Repo.DeleteNews(id); err != nil {
		logrus.WithError(err).Error("Failed to delete news from database")
		return err
	}

	// Hapus cache
	invalidateCache(ctx, cacheKey)
	logrus.WithField("news_id", id).Info("Cache invalidated after news deletion")
	return nil
}

// func (u *NewsUsecaseImpl) UpdateNews(ctx context.Context, news *entity.News) error {
// 	cacheKey := fmt.Sprintf("news:%d", news.ID)

// 	// Update data di database
// 	if err := u.Repo.UpdateNews(news); err != nil {
// 		logrus.WithError(err).Error("Failed to update news in database")
// 		return err
// 	}

// 	// Perbarui cache dengan data terbaru
// 	newsJSON, err := json.Marshal(news)
// 	if err != nil {
// 		logrus.WithError(err).Error("Failed to marshal news for cache")
// 		return err
// 	}

// 	if err := config.RedisClient.Set(ctx, cacheKey, newsJSON, 10*time.Second).Err(); err != nil {
// 		logrus.WithError(err).Error("Failed to update cache in Redis")
// 		return err
// 	}

// 	logrus.WithField("news_id", news.ID).Info("Cache updated after news update")
// 	return nil
// }
