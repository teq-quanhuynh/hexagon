package postgrestore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"hexagon/domain/book"
)

type BookStore struct {
	db *sqlx.DB
}

func NewBookStore(db *sqlx.DB) *BookStore {
	return &BookStore{db}
}

func (s *BookStore) Save(b *book.Book) error {
	_, err := s.db.Exec(`INSERT INTO books(isbn,name) VALUES ($1,$2)`, b.ISBN, b.Name)
	if err != nil {
		return fmt.Errorf("cannot save the book: %w", err)
	}
	return nil
}

func (s *BookStore) FindByISBN(isbn string) (*book.Book, error) {
	var result BookQuerySchema
	err := s.db.Get(&result, `SELECT isbn,name FROM books WHERE isbn=$1`, isbn)
	if err != nil {
		return nil, fmt.Errorf("cannot get the book '%s': %w", isbn, err)
	}

	b := book.NewBook(result.ISBN, result.Name)
	return &b, nil
}
