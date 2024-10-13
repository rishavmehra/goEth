package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ethClient *ethclient.Client

type AddressRequest struct {
	Address string `json:"address"`
}

func main() {
	var err error
	ethClient, err = EthCommon()
	if err != nil {
		log.Fatal("Error initializing Ethereum client: ", err)
	}

	r := gin.Default()
	r.GET("/", EthLatestBlock)
	r.POST("/balance", EthBalance)
	r.GET("/wallet", CreateWallet)

	r.Run()
}

func EthCommon() (*ethclient.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// infuraEndpoint := os.Getenv("INFURA_ENDPOINT")
	client, err := ethclient.Dial("/tmp/geth.ipc")
	if err != nil {
		log.Fatal("Unable to connect with the Client: ", err)
	}
	return client, nil
}

func CreateWallet(c *gin.Context) {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Private Key: ", privateKey)

	publicKey := privateKey.Public()
	fmt.Println("Public Key: ", publicKey)
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("New wallet Address:", address.Hex())
	c.JSON(http.StatusOK, gin.H{
		"New Wallet": address,
	})

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

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	c.JSON(http.StatusOK, gin.H{
		"balance": ethValue.String(), // need to convert tthis in string
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
