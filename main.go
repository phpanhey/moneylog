package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	options := &comdirect.AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	accountID := os.Getenv("COMDIRECT_ACCOUNT_ID")

	authenticator := comdirect.NewAuthenticator(options)
	client := comdirect.NewWithAuthenticator(authenticator)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := client.Authenticate(ctx)
	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	now := time.Now()
	firstOfMonth := time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		now.Location(),
	)

	opts := comdirect.EmptyOptions()
	opts.Add("min-bookingDate", firstOfMonth.Format("2006-01-02")).
		Add(comdirect.MaxBookingDateQueryKey, now.Format("2006-01-02")).
		Add(comdirect.PagingCountQueryKey, "100")

	txns, err := client.Transactions(ctx, accountID, opts)
	if err != nil {
		log.Fatalf("fetching transactions failed: %v", err)
	}

	for _, t := range txns.Values {
		fmt.Printf("%s  %8s %s  %s\n",
			t.BookingDate,
			t.Amount.Value,
			t.Amount.Unit,
			t.RemittanceInfo,
		)
	}
}
