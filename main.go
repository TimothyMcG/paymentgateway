package main

import (
	"fmt"
	"net/http"
	"paymentgateway/bank"
	"paymentgateway/db"
	"paymentgateway/model"
	"paymentgateway/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

var database db.Db

func main() {

	conn := db.SetUpDB()
	database = conn

	router := SetUpRouter()
	router.Run()
}

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")
	public.POST("/register", Register)
	public.POST("/login", Login)

	protected := r.Group("/api/v1")
	protected.Use(jwtAuth())
	protected.POST("/payment", SendPayment)
	protected.GET("/payment", GetPayment)
	return r
}

func SendPayment(ctx *gin.Context) {
	var payment model.Payment

	err := ctx.BindJSON(&payment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := payment.Validate()
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment"})
		return
	}

	req := bank.NewBankRequest()
	req.Amount = payment.Amount
	req.CardName = payment.CardName
	req.CardNumber = payment.CardNumber
	req.CardType = payment.CardType
	req.Currency = payment.Currency
	req.Description = payment.Description
	id, success := req.ProcessPayment()
	if !success {
		ctx.JSON(http.StatusPaymentRequired, gin.H{"error": "payment failed"})
		return
	}

	payment.Id = id
	err = database.SavePayment(payment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save payment"})
		return
	}

	ctx.JSON(200, gin.H{"paymentID": payment.Id})
}

func GetPayment(ctx *gin.Context) {

	idString := ctx.Query("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	payment, err := database.GetPayment(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	payment.Masker()

	ctx.JSON(200, payment)

}

func jwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := util.TokenValid(ctx)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func Register(ctx *gin.Context) {
	var merchant model.Merchant
	err := ctx.BindJSON(&merchant)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, found, err := database.VerifyMerchant(merchant)
	if found {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "merchant already exists"})
		return
	}

	err = database.SaveUser(merchant)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"merchant": merchant.Name})

}

func Login(ctx *gin.Context) {
	var merchant model.Merchant
	err := ctx.BindJSON(&merchant)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, match, err := database.VerifyMerchant(merchant)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if !match {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "merchant not recognised"})
		return
	}

	token, err := util.GenerateToken(uint(merchant.Id))
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "merchantID": id})
}
