package usecase

import (
	"context"

	"github.com/kusnadi8605/news/entity"
)

type NewsUsecase interface {
	GetAllNews(ctx context.Context) (*[]entity.News, error)
	GetNewsByID(ctx context.Context, id int) (*entity.News, error)
	CreateNews(ctx context.Context, news *entity.News) (int64, error)
	UpdateNews(ctx context.Context, news *entity.News) error
	DeleteNews(ctx context.Context, id int) error
}
