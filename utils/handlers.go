package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StoreHandler(c *gin.Context) {

}

// Health ...
func Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
