package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"netdisk/model"
	"netdisk/route"
)

func main() {
	db, err := model.InitDB()
	if err != nil {
		log.Println("err open databases", err)
		return
	}
	defer db.Close()

	gin := gin.New()

	route.LOAD(
		gin,
	)

	gin.Run()
}