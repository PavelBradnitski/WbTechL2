package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const fileName = "requests.log"

// LoggerToFile — middleware для логирования в файл
func LoggerToFile() gin.HandlerFunc {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ошибка открытия файла для логов: %v", err)
		return func(c *gin.Context) {
			c.Next() // Continue processing the request even if logging fails
		}
	}

	logger := log.New(file, "", log.LstdFlags)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Printf("%s %s | %v", c.Request.Method, c.Request.URL.Path, duration)
	}
}
