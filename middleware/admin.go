package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil header auth
		key := os.Getenv("ADMIN_SECRET")

		// TODO: validasi header dengan kata sandi admin
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Akses tidak di izinkan"})
			c.Abort()
			return
		}

		if auth != key {
			c.JSON(401, gin.H{"error": "Akses tidak di izinkan"})
			c.Abort()
			return
		}

		// TODO: lanjutkan request ke handler
		c.Next()
	}
}
