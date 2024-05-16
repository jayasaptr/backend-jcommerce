package model

import (
	"database/sql"
	"errors"
)

type Products struct {
	ID        string `json:"id" binding:"len=0"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

var (
	ErrDBNil = errors.New("koneksi tidak tersedia")
)

func SelectProduct(db *sql.DB) ([]Products, error) {
	if db == nil {
		return nil, ErrDBNil
	}
	query := `SELECT id, name, price FROM products WHERE is_deleted = false;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []Products{}
	for rows.Next() {
		var product Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func SelectProductByID(db *sql.DB, id string) (Products, error) {
	if db == nil {
		return Products{}, ErrDBNil
	}

	query := `SELECT id, name, price FROM products WHERE is_deleted = false AND id = $1`

	var product Products
	row := db.QueryRow(query, id)
	if err := row.Scan(&product.ID, &product.Name, &product.Price); err != nil {
		return Products{}, err
	}

	return product, nil
}

func InsertProduct(db *sql.DB, product Products) error {
	if db == nil {
		return ErrDBNil
	}

	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3);`
	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(db *sql.DB, product Products) error {
	if db == nil {
		return ErrDBNil
	}

	query := `UPDATE products SET name=$1, price=$2 WHERE id=$3;`
	_, err := db.Exec(query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}
