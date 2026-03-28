package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/afrochainorg/afrochain-faucet/config"
	"github.com/afrochainorg/afrochain-faucet/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var config config.Config

	flag.StringVar(&config.FaucetPort, "port", "", "example: 9000")

	flag.StringVar(&config.BinaryName, "cli", "", "example: gaiad, neutrond, etc....")
	flag.StringVar(&config.FaucetWalletAddress, "address", "", "example: cosmos1zypqa76j...")
	flag.StringVar(&config.FaucetWalletAlias, "alias", "alice", "example: alice")
	flag.StringVar(&config.ChainNode, "node", "", "example: http://localhost:27657")
	flag.StringVar(&config.ChainHome, "home", "", "example: ~/.gaiad")
	flag.StringVar(&config.KeyringBackend, "keyring-backend", "test", "example: test")
	flag.StringVar(&config.ChainID, "chain-id", "", "example: cosmoshub-test-1")

	flag.Parse()

	if err := config.IsValid(); err != nil {
		log.Fatalln(err)
	}

	h, err := handler.New(config)
	if err != nil {
		log.Fatalln("Failed to initialize handler:", err)
	}

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/request", h.Request)
	r.Run(fmt.Sprintf(":%s", config.FaucetPort))
}
