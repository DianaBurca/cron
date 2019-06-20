package main

import (
	"./utils"
	"github.com/gin-gonic/gin"
)

func main() {

	driver := gin.Default()

	driver.PUT("/store/:name", utils.StoreHandler)

	driver.Run()

}
