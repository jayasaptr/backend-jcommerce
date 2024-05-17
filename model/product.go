package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func DeleteProduct(db *sql.DB, id string) error {
	if db == nil {
		return ErrDBNil
	}

	query := `UPDATE products SET is_deleted = true WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func SelectProductIn(db *sql.DB, ids []string) ([]Products, error) {
	if db == nil {
		return nil, ErrDBNil
	}

	if len(ids) == 0 {
		return []Products{}, nil
	}

	placeholders := []string{}
	arg := []interface{}{}

	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		arg = append(arg, id)
	}

	query := fmt.Sprintf(`SELECT id, name, price FROM products WHERE is_deleted = false AND id IN (%s);`, strings.Join(placeholders, ","))
	rows, err := db.Query(query, arg...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Products{}
	for rows.Next() {
		var product Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
