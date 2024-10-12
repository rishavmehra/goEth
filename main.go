package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ethClient *ethclient.Client

func main() {
	var err error
	ethClient, err = EthCommon()
	if err != nil {
		log.Fatal("Error initializing Ethereum client: ", err)
	}

	r := gin.Default()
	r.GET("/", EthLatestBlock)
	r.POST("/balance", EthBalance)

	r.Run()
}

type AddressRequest struct {
	Address string `json:"address"`
}

func EthBalance(c *gin.Context) {
	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	address := common.HexToAddress(req.Address)

	balance, err := ethClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal("Unable to fetch balance: ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance.String(), // need to convert tthis in string
	})
}

func EthLatestBlock(c *gin.Context) {
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Println("Error while getting the block number:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get block number",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message for block number": blockNumber,
	})
}

func EthCommon() (*ethclient.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infuraEndpoint := os.Getenv("INFURA_ENDPOINT")
	client, err := ethclient.Dial(infuraEndpoint)
	if err != nil {
		log.Fatal("Unable to connect with the Client: ", err)
	}
	return client, nil
}
