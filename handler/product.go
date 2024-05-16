package handler

import (
	"database/sql"
	"errors"
	"log"
	"toko-online/model"

	"github.com/gin-gonic/gin"
)

func ListProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil dari database berikan response
		products, err := model.SelectProduct(db)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil data product %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}
		// TODO: berikan response
		c.JSON(200, products)

	}

}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: baca id dari url
		id := c.Param("id")
		// TODO: ambil dari database dengan id
		product, err := model.SelectProductByID(db, id)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Terjadi kesalahan saat mengambil data product %v\n", err)
				c.JSON(404, gin.H{"error": "Product tidak ditemukan"})
				return
			}

			log.Printf("Terjadi kesalahan saat mengambil data product %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// TODO: berikan response
		c.JSON(200, product)
	}

}
