package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StorePayload ...
type StorePayload struct {
	CityName string `json:"city_name"`
}

// StoreHandler ...
func StoreHandler(c *gin.Context) {
	payload := StorePayload{}
	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		fmt.Println(err)
		return
	}
	status := findOrCreate(payload.CityName)

	c.Status(status)

}

// Health ...
func Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
