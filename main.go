package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
	// "log"
)

type Account struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Balance int
}

// Tracks the randomnly generated account IDs
var acctIDs []uuid.UUID

// Insert new rows into the accounts table
// This function generates new UUIDs and random balances for each row, and
// then it appends the ID to the `acctIDs`, which other functions use to track the IDs
func addAccounts(db *gorm.DB, numRows int, transferAmount int) error {
	log.Printf("Creating %d new accounts...", numRows)

	for i := 0; i < numRows; i++ {
		newId := uuid.New()
		newBalance := rand.Intn(1000) + transferAmount
		if err := db.Create(&Account{ID: newId, Balance: newBalance}).Error; err != nil {
			return err
		}
		acctIDs = append(acctIDs, newId)
	}
	log.Println("Accounts created")
	return nil
}

// Transfer funds between accounts
// This function adds `amount` to the `balance` column of the row with the "id" column matching `toId`,
// and removes `amount` from the `balance` column of the row with the "id" column matching `fromId`.
func transferFunds(db *gorm.DB, fromId uuid.UUID, toId uuid.UUID, amount int) error {
	log.Printf("Transferring %d from account %s to account %s...", amount, fromId, toId)
	var fromAccount Account
	var toAccount Account

	db.First(&fromAccount, fromId)
	db.First(&toAccount, toId)

	if fromAccount.Balance < amount {
		return fmt.Errorf("account %s balance %d is lower than transfer amount %d", fromAccount.ID, fromAccount.Balance, amount)
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if err := db.Save(&fromAccount).Error; err != nil {
		return err
	}

	if err := db.Save(&toAccount).Error; err != nil {
		return err
	}

	log.Println("Funds transfered")
	return nil
}

func printBalance(db *gorm.DB) {
	var accounts []Account
	db.Find(&accounts)
	fmt.Printf("Balance at '%s:\n", time.Now())
	for _, account := range accounts {
		fmt.Printf("%s %d\n", account.ID, account.Balance)
	}
}

// Delete all rows in "accounts" table inserted by 'main': (i.e., tracked by `acctIDs`)
func deleteAccounts(db *gorm.DB, accountIDs []uuid.UUID) error {
	log.Println("Delete accounts created...")
	err := db.Where("id IN ?", accountIDs).Delete(Account{}).Error
	if err != nil {
		return err
	}
	log.Println("Account deleted.")
	return nil
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	fmt.Println("Hello this is a rest api")
	fmt.Printf("DATABASE_URL: %s\n", os.Getenv("DATABASE_URL"))
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")+"&application_name=docs_simplecrud_gorm"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// Automatically create the "account" table based on the Account model
	db.AutoMigrate(&Account{})

	// The number of inital rows to insert

	const numAccts int = 5

	// The amount to transfer between accounts
	const transferAmt int = 100

	// Insert `numAccts` rows into the "accounts" table.
	// To handle potential transaction retries, we wrap the call to `addAcounts` in `crdbgorm.ExecuteTx`, a helper function for GORM which implements a retry loop
	// err = crdbgorm.ExecuteTx(context.Background(), db, func(tx *gorm.DB) error {
	// 	return addAccounts(tx, numAccts, transferAmt)
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err := crdbgorm.ExecuteTx(context.Background(), db, nil,
		func(tx *gorm.DB) error {
			return addAccounts(db, numAccts, transferAmt)
		},
	); err != nil {
		fmt.Println(err)
	}

	// Print balances before transfer
	printBalance(db)

	// Selecct two accoutn Ids
	fromID := acctIDs[0]
	toID := acctIDs[0:][rand.Intn(len(acctIDs))]

	// Transfer funds between accounts. To handle potential
	// transaction retry errors, we wrap the call to 'trnasferFunds' in `crdbgorm.ExecuteTx`

	if err := crdbgorm.ExecuteTx(context.Background(), db, nil,
		func(tx *gorm.DB) error {
			return transferFunds(tx, fromID, toID, transferAmt)
		},
	); err != nil {
		fmt.Println(err)
	}

	// Print balances after transfer to ensure that it worked
	printBalance(db)

	// Delete all accounts created by the earlier call to `addAccounts`
	// To handle potential transaction retry errors, we wrap the call
	// to `deleteAccounts` in `crdbgorm.ExecuteTx`
	if err := crdbgorm.ExecuteTx(context.Background(), db, nil,
		func(tx *gorm.DB) error {
			return deleteAccounts(db, acctIDs)
		},
	); err != nil {
		fmt.Println(err)
	}
}
