package handler

import (
	"database/sql"
	"errors"
	"log"
	"toko-online/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CraeteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Products
		if err := c.Bind(&product); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data Product tidak valid"})
			return
		}

		product.ID = uuid.New().String()

		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Terjadi kesalahan pada server saat create product: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var productReq model.Products
		if err := c.Bind(&productReq); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data Product tidak valid"})
			return
		}

		product, err := model.SelectProductByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil product: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if productReq.Name != "" {
			product.Name = productReq.Name
		}

		if productReq.Price != 0 {
			product.Price = productReq.Price
		}

		if err := model.UpdateProduct(db, product); err != nil {
			log.Printf("Terjadi kesalahan saat mengupdate product: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
