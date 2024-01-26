package main

import (
	"errors"
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"time"
)

type Card struct {
	gorm.Model
	Number  string
	Pin     string
	Balance int `gorm:"default:0"`
}

func randomIntArr(slice []byte) {
	for i := range slice {
		slice[i] = byte(48 + rand.Intn(10))
	}
}

// generating checksum using Luhn algorithm
func generateChecksum(first15Digits []byte) int {
	copySlice := make([]byte, len(first15Digits))
	copy(copySlice, first15Digits)
	for i := range copySlice {
		copySlice[i] -= 48
		if (i+1)%2 == 1 {
			copySlice[i] *= 2
		}
		if copySlice[i] > 9 {
			copySlice[i] -= 9
		}
	}

	total := 0
	for _, el := range copySlice {
		total += int(el)
	}

	checksum := 10 - (total % 10)
	if checksum == 10 {
		return 0
	}
	return checksum
}

func verifyCard(cardNum string) bool {
	lastIndex := len(cardNum) - 1
	first15 := []byte(cardNum[:lastIndex])
	cardChecksum := int(cardNum[lastIndex] - 48)
	if cardChecksum == generateChecksum(first15) {
		return true
	} else {
		return false
	}
}

func createAccount() (string, string) {
	bin := "400000"
	accountNum := make([]byte, 9)
	pin := make([]byte, 4)

	randomIntArr(accountNum)
	randomIntArr(pin)

	first15 := append([]byte(bin), accountNum...)
	checksum := generateChecksum(first15)
	creditCardNum := fmt.Sprintf("%s%d", string(first15), checksum)

	fmt.Println("\nYour card has been created")
	fmt.Println("Your card number:\n" + creditCardNum)
	fmt.Printf("Your card pin:\n%s", string(pin))
	fmt.Println("\n")

	return creditCardNum, string(pin)
}

func logIn(db *gorm.DB) (bool, uint) {
	var inputCardNum string
	var inputPin string

	fmt.Println("\nEnter your card number:")
	fmt.Scan(&inputCardNum)
	fmt.Println("Enter your PIN:")
	fmt.Scan(&inputPin)

	var card Card
	result := db.Where("number = ?", inputCardNum).Where("pin = ?", inputPin).First(&card)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("\nWrong card number or PIN!\n")
		return false, 0
	} else {
		fmt.Println("\nYou have successfully logged in!\n")
		return true, card.ID
	}
}

func handleTransfer(receiverCardNum string, sender *Card, db *gorm.DB) {
	// if transferring to same card, then abort
	if receiverCardNum == sender.Number {
		fmt.Println("You can't transfer money to the same account!\n")
		return
	}

	// retrieving the receiver card details
	var receiver Card
	result := db.Where("number = ?", receiverCardNum).First(&receiver)
	// if transferring to card that doesn't exist in the database, then abort
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Such a card does not exist.\n")
		return
	}

	var amount int
	fmt.Println("Enter how much money you want to transfer:")
	fmt.Scan(&amount)
	// if transferring more money than the balance, then abort
	if amount > sender.Balance {
		fmt.Println("Not enough money!\n")
		return
	}

	// performing transfer using transaction
	tx := db.Begin()

	sender.Balance -= amount
	result = tx.Save(sender)
	if result.Error != nil {
		log.Printf("cannot update Sender's balance: %v\n", result.Error)
		tx.Rollback()
		return
	}

	receiver.Balance += amount
	result = tx.Save(&receiver)
	if result.Error != nil {
		log.Printf("cannot update Receiver's balance: %v\n", result.Error)
		tx.Rollback()
		return
	}
	fmt.Println("Success!\n")

	tx.Commit()
}

func appMenu(db *gorm.DB, cardID uint) {
	var card Card
	result := db.Where("id = ?", cardID).First(&card)
	if result.Error != nil {
		log.Fatalf("cannot find Card: %v\n", result.Error)
	}

	var userInput int
	for {
		fmt.Println("1. Balance")
		fmt.Println("2. Add income")
		fmt.Println("3. Do transfer")
		fmt.Println("4. Close account")
		fmt.Println("5. Log out")
		fmt.Println("0. Exit")

		fmt.Scan(&userInput)
		switch userInput {
		case 0:
			// Exit
			os.Exit(0)
		case 1:
			// Balance
			fmt.Printf("\nBalance: %d\n\n", card.Balance)
		case 2:
			// Add income
			var income int
			fmt.Println("\nEnter income:")
			fmt.Scan(&income)
			card.Balance += income
			result := db.Save(&card)
			if result.Error != nil {
				log.Fatalf("cannot update Card: %v\n", result.Error)
			}
			fmt.Println("Income was added!\n")
		case 3:
			// Do transfer
			var receiverCardNum string
			fmt.Println("\nTransfer")
			fmt.Println("Enter card number:")
			fmt.Scan(&receiverCardNum)
			if verifyCard(receiverCardNum) {
				handleTransfer(receiverCardNum, &card, db)
			} else {
				fmt.Println("Probably you made a mistake in the card number. Please try again!\n")
			}
		}

		// Close account
		if userInput == 4 {
			result := db.Delete(&card)
			if result.Error != nil {
				log.Fatalf("cannot delete Card: %v\n", result.Error)
			}
			fmt.Println("\nThe account has been closed!\n")
			break
		}

		// Log out
		if userInput == 5 {
			fmt.Println("\nYou have successfully logged out!\n")
			break
		}
	}
}

func main() {
	// DO NOT delete the `rand.Seed(...)` line, it initializes the random number generator!
	rand.Seed(time.Now().UnixNano())

	// Getting the database filename from terminal
	dbName := flag.String("fileName", "example.db", "Type the database name.")
	flag.Parse()
	// Creating the database schema
	db, err := gorm.Open(sqlite.Open(*dbName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Migrating the Card struct to create the cards table in the database
	err = db.AutoMigrate(&Card{})
	if err != nil {
		log.Fatal(err)
	}

	// App
	var userInput int
	for {
		fmt.Println("1. Create an account")
		fmt.Println("2. Log into account")
		fmt.Println("0. Exit")

		fmt.Scan(&userInput)
		switch userInput {
		case 0:
			// Exit
			os.Exit(0)
		case 1:
			// Create card
			creditCardNum, pin := createAccount()
			card := Card{Number: creditCardNum, Pin: pin}
			result := db.Create(&card)
			if result.Error != nil {
				log.Fatalf("cannot create Card: %v\n", result.Error)
			}
		case 2:
			// Login
			success, cardID := logIn(db)
			if success {
				appMenu(db, cardID)
			}
		}
	}
}
