package main

import (
	"fmt"
	"os"
)


func calcTotalIncome(rows, seats int) int {
	var income int
	totalSeats := rows * seats

	if totalSeats > 60 {
		if rows % 2 == 1 {
			frontHalf := rows / 2
			backHalf := frontHalf + 1
			income = frontHalf * seats * 10 + backHalf * seats * 8
		} else {
			half := rows / 2
			income = half * seats * 10 + half * seats * 8
		}
	} else {
		income = totalSeats * 10
	}

	return income
}


func showSeating(seating [][]string, seats int) {
	fmt.Println("\nCinema:")

	firstStr := " "
	for i := 1; i <= seats; i++ {
		temp := fmt.Sprintf(" %d", i)
		firstStr += temp
	}
	fmt.Println(firstStr)

	for i, el := range seating {
		rowStr := ""
		temp := fmt.Sprintf("%d", i + 1)
		rowStr += temp
		for _, seat := range el {
			rowStr += " "
			rowStr += seat
		}
		fmt.Println(rowStr)
	}
}


func bookTicket(rows int, seats int, seating [][]string, ticketsSold *int, currentIncome *int) {
	var rowNum, seatNum, ticketPrice int
	totalSeats := rows * seats

	// taking input and checking seat availability
	for {
		fmt.Print("\nEnter a row number: ")
		fmt.Scan(&rowNum)
		fmt.Print("Enter a seat number in that row: ")
		fmt.Scan(&seatNum)

		if (rowNum < 1 || rowNum > rows) || (seatNum < 1 || seatNum > seats) {
			// if the seat is out of bounds
			fmt.Println("\nWrong input!")
			
		} else if seating[rowNum - 1][seatNum - 1] == "B" {
			// if the seat is already booked
			fmt.Println("\nThat ticket has already been purchased!")

		} else {
			break
		}
	}


	// calculating ticket price based on position in cinema
	if totalSeats > 60 {
		front := (rows / 2) * seats
		pos := (rowNum - 1) * seats + seatNum

		if pos <= front {
			ticketPrice = 10
		} else {
			ticketPrice = 8
		}
	} else {
		ticketPrice = 10
	}

	fmt.Printf("Ticket price: $%d\n", ticketPrice)
	*ticketsSold += 1
	*currentIncome += ticketPrice


	// updating and re-printing the seatings
	seating[rowNum - 1][seatNum - 1] = "B"
}


func main() {
	var rows, seats, ticketsSold, currentIncome int

	// taking input
    fmt.Print("Enter the number of rows: ")
	fmt.Scan(&rows)
    fmt.Print("Enter the number of seats in each row: ")
	fmt.Scan(&seats)

	totalIncome := calcTotalIncome(rows, seats)


	// constructing slice which will store the seating data
	var seating [][]string = make([][]string, rows)

	for i := range seating {
		seating[i] = make([]string, seats)
		for j := range seating[i] {
			seating[i][j] = "S"
		}
	}


	// MENU
	for {
		var userInput int

		fmt.Println("\n1. Show the seats")
		fmt.Println("2. Buy a ticket")
		fmt.Println("3. Statistics")
		fmt.Println("0. Exit")
		fmt.Scan(&userInput)

		switch userInput {
		case 1:
			showSeating(seating, seats)
		case 2:
			bookTicket(rows, seats, seating, &ticketsSold, &currentIncome)
		case 3:
			fmt.Println("\nNumber of purchased tickets:", ticketsSold)

			percentage := (float32(ticketsSold) / float32(rows * seats)) * 100
			fmt.Printf("Percentage: %.2f%%\n", percentage)

			fmt.Printf("Current income: $%d\n", currentIncome)
			fmt.Printf("Total income: $%d\n", totalIncome)
		case 0:
			os.Exit(0)
		}
	}
}