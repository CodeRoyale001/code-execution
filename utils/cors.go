package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() func(c *gin.Context) {
	allowedOrigins := []string{"http://localhost:3000", "https://serene-fortress-91389-77d1fb95872a.herokuapp.com", "https://oj-front-end.vercel.app"}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is present and allowed
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, origin, x-requested-with")
		// Add AllowCredentials: true
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
