package main

import (
	"Signer/internal/db"
	"Signer/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()
	r := gin.Default()
	r.POST("/register", handlers.RegisterUser)
	r.POST("/sign_answers", handlers.SignAnswers)
	r.GET("/verify_signature", handlers.VerifySignature)

	r.Run()

}
