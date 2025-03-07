package repository

import "github.com/kusnadi8605/news/entity"

type NewsRepository interface {
	GetAllNews() (*[]entity.News, error)
	GetNewsByID(id int) (*entity.News, error)
	CreateNews(news *entity.News) (int64, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id int) error
}
