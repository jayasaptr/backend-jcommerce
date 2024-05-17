package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
	"toko-online/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: mengambil data pesanan dari request body
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			log.Printf("Terjadi kesalahan ketika membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Terjadi kesalahan pada request"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)
		for _, o := range checkoutOrder.Products {
			ids = append(ids, o.ID)
			orderQty[o.ID] = int32(o.Quantity)
		}

		// TOOD: mengambil data dari db
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil product: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
		}

		// TODO: membuat kata sandi
		passcode := generatedPasscode(5)

		// TODO: hash kata sandi
		hashCode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membuat hash: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		hashCodeString := string(hashCode)

		// TODO: buat order dan detail

		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			Passcode:   &hashCodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		for _, p := range products {
			total := p.Price * int64(orderQty[p.ID])

			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     total,
			}

			details = append(details, detail)

			order.GrandTotal += total
		}

		model.CraeteOrder(db, order, details)

		order.Passcode = &passcode

		ordeWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}

		c.JSON(200, ordeWithDetail)
	}
}

func generatedPasscode(lenght int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, lenght)

	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil id
		id := c.Param("id")

		// TODO: baca req body
		var confirmRequest model.Confrim
		if err := c.BindJSON(&confirmRequest); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body %v\n", err)
			c.JSON(400, gin.H{"error": "Data pesanan tidak valid"})
			return
		}

		// TODO: ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca data pesanan %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}
		// TODO: cocokan kata sandi
		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirmRequest.Passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat mecocokan kata sandi %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak di izinkan mengakses pesanan"})
			return
		}

		// TODO: pastikan pesanan belum di bayar
		if order.PaidAt != nil {
			log.Println("Pesanan sudah dibayar")
			c.JSON(400, gin.H{"error": "Pesanan sudah di bayar"})
			return
		}

		// TODO: cocokan jumlah pembayaran
		if order.GrandTotal != confirmRequest.Amount {
			log.Printf("Jumlah harga tidak sesuai %v\n", err)
			c.JSON(400, gin.H{"error": "Jumlah pembayaran tidak sesuai"})
			return
		}

		// TODO: update informasi pesanan
		current := time.Now()

		if err = model.UpdateOrderByID(db, id, confirmRequest, current); err != nil {
			log.Printf("Terjadi kesalahan saat mengupdate data pesanan %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		order.Passcode = nil

		order.PaidAt = &current
		order.PaidBank = &confirmRequest.Bank
		order.PaidAccountNumber = &confirmRequest.AccountNumber

		c.JSON(200, order)
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil id
		id := c.Param("id")

		// TODO: ambil passcode dari query param
		passcode := c.Query("passcode")

		// TODO: ambil data dari databse
		order, err := model.SelectOrderByID(db, id)

		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca data pesanan %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode tidak valid")
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
		}

		log.Println(order.Passcode)

		// TODO: cocokan kata sandi pesanan
		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat mecocokan kata sandi %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak di izinkan mengakses pesanan"})
			return
		}

		order.Passcode = nil
		c.JSON(200, order)
	}
}
