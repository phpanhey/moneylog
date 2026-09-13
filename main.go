package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
)

func main() {

	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	options := &comdirect.AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	authenticator := comdirect.NewAuthenticator(options)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := authenticator.Authenticate(ctx)

	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	client := comdirect.NewWithAuthenticator(authenticator)
	_ = client

}
