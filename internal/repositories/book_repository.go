package repositories

import (
	"gorm-test/internal/models"

	"gorm.io/gorm"
)

//go:generate mockery --name BookWriter
type BookWriter interface {
	Insert(book *models.Book) error
}

type BookRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		DB: db,
	}
}

func (t *BookRepository) Insert(book *models.Book) error {
	return t.DB.Create(book).Error
}

func (t *BookRepository) Get(id uint) (*models.Book, error) {
	var book *models.Book

	err := t.DB.First(&book, id).Error

	if err != err {
		return nil, err
	}

	return book, nil
}
