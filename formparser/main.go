package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	// Open the input CSV file
	inputFile, err := os.Open("sign_up_form_2024-10-01.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// Create a reader for the CSV
	reader := csv.NewReader(inputFile)

	// Read all rows from the input CSV
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Open the output CSV file
	fileName := fmt.Sprintf("output-%s.csv", time.Now().Format("2006-01-02-1504"))
	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Create a writer for the output CSV
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Get headers from the input CSV (first row)
	headers := rows[0]

	// removed column prefix "adult1_"
	for i, colName := range headers {
		headers[i] = strings.ReplaceAll(colName, "adult1_", "")
	}

	// Create header
	adultHeader := []string{}
	adultHeader = append(adultHeader, "adult_id")        // generated index for the adult
	adultHeader = append(adultHeader, "household_id")    // generated index for the household
	adultHeader = append(adultHeader, headers[53])       // adult name
	adultHeader = append(adultHeader, "email")           // adult email
	adultHeader = append(adultHeader, headers[54:70]...) // survey fields
	adultHeader = append(adultHeader, "anything_else")   // other commend
	adultHeader = moveColumns16to18After5(adultHeader)
	err = writer.Write(adultHeader)
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}

	// Process each row
	personIndex := 0
	householdIndex := 0
	for _, row := range rows[1:] {
		// Write columns 10 to 14 to the new CSV
		personIndex++
		householdIndex++
		adultRow := []string{}
		adultRow = append(adultRow, strconv.Itoa(personIndex))
		adultRow = append(adultRow, strconv.Itoa(householdIndex))
		adultRow = append(adultRow, row[53])       // adult1 name
		adultRow = append(adultRow, row[1])        // email from beginning of form
		adultRow = append(adultRow, row[54:70]...) // remaining adult1 field
		adultRow = append(adultRow, row[89:90]...) // duplicate of 'anything else' since only one per household
		adultRow[4] = replaceParticipationLevel(adultRow[4])
		adultRow = moveColumns16to18After5(adultRow)
		err = writer.Write(adultRow)
		if err != nil {
			fmt.Println("Error writing row:", err)
			return
		}

		if row[70] == "Yes" {
			personIndex++
			adultRow = []string{}
			adultRow = append(adultRow, strconv.Itoa(personIndex))
			adultRow = append(adultRow, strconv.Itoa(householdIndex))
			adultRow = append(adultRow, row[71:90]...)
			adultRow[4] = replaceParticipationLevel(adultRow[4])
			adultRow = moveColumns16to18After5(adultRow)
			err = writer.Write(adultRow)
			if err != nil {
				fmt.Println("Error writing additional row:", err)
				return
			}
		}
	}

	fmt.Println("CSV processing completed successfully.")
}

// moveColumns15to17After4 takes a slice, cuts the columns 16, 17, 18 (index 15, 16, 17)
// and inserts them after the 4th element (index 3).  This moves the availability columns
func moveColumns16to18After5(row []string) []string {
	cutFrom := 16
	cutTo := 18
	insertInto := 5
	if len(row) < cutTo+1 {
		// If the slice has fewer than 19 elements, return the original slice unchanged.
		return row
	}

	// Cut elements
	cols20to22 := row[cutFrom : cutTo+1]

	// Keep the part of the slice before the cut elements
	// And also the part after the cut elements
	sliceBefore := row[:cutFrom]
	sliceAfter := row[cutTo+1:]

	result := sliceBefore[:insertInto]
	result = append(result, cols20to22...)
	result = append(result, sliceBefore[:insertInto]...)
	result = append(result, sliceAfter...)

	return result
}

// Helper function to replace text in the 3rd column
func replaceParticipationLevel(col string) string {
	if strings.Contains(col, "I want to lead a class") {
		return "Can lead"
	} else if strings.Contains(col, "want to share responsibility") {
		return "Can lead with support"
	} else if strings.Contains(col, "I want to help support") {
		return "Can help"
	} else if strings.Contains(col, "I am not available for any") {
		return "Not available in fall"
	} else if strings.Contains(col, "I have an extenuating circumstance") {
		return "Not available in fall"
	}
	return col
}
