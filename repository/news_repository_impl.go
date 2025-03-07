package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kusnadi8605/news/entity"
)

type NewsRepositoryImpl struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) NewsRepository {
	return &NewsRepositoryImpl{db: db}
}

func (r *NewsRepositoryImpl) CreateNews(news *entity.News) (int64, error) {
	result, err := r.db.Exec("INSERT INTO news (title, content, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		news.Title, news.Content)
	if err != nil {
		return 0, err
	}

	// get id dari data yang baru di input
	newsID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return newsID, nil
}

func (r *NewsRepositoryImpl) UpdateNews(news *entity.News) error {
	_, err := r.db.Exec("UPDATE news SET title = ?, content = ?, updated_at = NOW() WHERE id = ?", news.Title, news.Content, news.ID)
	return err
}

func (r *NewsRepositoryImpl) DeleteNews(id int) error {
	result, err := r.db.Exec("DELETE FROM news WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		err := fmt.Errorf("news with ID %d not found", id)
		return err
	}

	return nil
}

func (r *NewsRepositoryImpl) GetNewsByID(id int) (*entity.News, error) {
	var news entity.News
	var createdAtStr, updatedAtStr string

	err := r.db.QueryRow("SELECT id, title, content, created_at, updated_at FROM news WHERE id = ?", id).
		Scan(&news.ID, &news.Title, &news.Content, &createdAtStr, &updatedAtStr)

	if err != nil {
		return nil, err
	}

	// Konversi string timestamp ke time.Time
	news.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	news.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

	return &news, nil
}

func (r *NewsRepositoryImpl) GetAllNews() (*[]entity.News, error) {
	rows, err := r.db.Query("SELECT id, title, content, created_at, updated_at FROM news")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []entity.News
	for rows.Next() {
		var news entity.News
		var createdAtStr, updatedAtStr string

		// Scan data
		if err := rows.Scan(&news.ID, &news.Title, &news.Content, &createdAtStr, &updatedAtStr); err != nil {
			return nil, err
		}

		// Konversi string ke time.Time
		news.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		news.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		newsList = append(newsList, news)
	}

	return &newsList, nil
}
