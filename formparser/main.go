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

	// Open the output CSV files
	timestamp := time.Now().Format("2006-01-02-1504")
	fileName := fmt.Sprintf("adults-%s.csv", timestamp)
	adultsFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating adults file:", err)
		return
	}
	defer adultsFile.Close()

	fileName = fmt.Sprintf("students-%s.csv", timestamp)
	studentsFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating students file:", err)
		return
	}
	defer adultsFile.Close()

	// Create writers for the output CSVs
	adultsFileWriter := csv.NewWriter(adultsFile)
	defer adultsFileWriter.Flush()

	studentsFileWriter := csv.NewWriter(studentsFile)
	defer studentsFileWriter.Flush()

	// Get headers from the input CSV (first row)
	headers := rows[0]

	// remove column prefix "adult1_"
	for i, colName := range headers {
		headers[i] = strings.ReplaceAll(colName, "adult1_", "")
	}

	// remove column prefix "child1_"
	for i, colName := range headers {
		headers[i] = strings.ReplaceAll(colName, "child1_", "")
	}

	// Create header
	adultsHeader := []string{}
	adultsHeader = append(adultsHeader, "adult_id")        // generated index for the adult
	adultsHeader = append(adultsHeader, "household_id")    // generated index for the household
	adultsHeader = append(adultsHeader, headers[53])       // adult name
	adultsHeader = append(adultsHeader, "email")           // adult email
	adultsHeader = append(adultsHeader, headers[54:70]...) // survey fields
	adultsHeader = append(adultsHeader, "anything_else")   // other commend
	adultsHeader = moveColumns16to18After5(adultsHeader)
	err = adultsFileWriter.Write(adultsHeader)
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}
	// Create header
	studentsHeader := []string{}
	studentsHeader = append(studentsHeader, "student_id")     // generated index for the adult
	studentsHeader = append(studentsHeader, "household_id")   // generated index for the household
	studentsHeader = append(studentsHeader, headers[2])       // student full name
	studentsHeader = append(studentsHeader, headers[3:14]...) // interests
	err = studentsFileWriter.Write(studentsHeader)
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}

	// Process each row
	adultIndex := 0
	studentIndex := 0
	householdIndex := 0
	for _, row := range rows[1:] {
		// Write adult fields
		adultIndex++
		householdIndex++
		adultRow := []string{}
		adultRow = append(adultRow, strconv.Itoa(adultIndex))
		adultRow = append(adultRow, strconv.Itoa(householdIndex))
		adultRow = append(adultRow, row[53])       // adult1 name
		adultRow = append(adultRow, row[1])        // email from beginning of form
		adultRow = append(adultRow, row[54:70]...) // remaining adult1 field
		adultRow = append(adultRow, row[89:90]...) // duplicate of 'anything else' since only one per household
		adultRow[4] = replaceParticipationLevel(adultRow[4])
		adultRow = moveColumns16to18After5(adultRow)
		err = adultsFileWriter.Write(adultRow)
		if err != nil {
			fmt.Println("Error writing row:", err)
			return
		}

		if row[70] == "Yes" {
			adultIndex++
			adultRow = []string{}
			adultRow = append(adultRow, strconv.Itoa(adultIndex))
			adultRow = append(adultRow, strconv.Itoa(householdIndex))
			adultRow = append(adultRow, row[71:90]...)
			adultRow[4] = replaceParticipationLevel(adultRow[4])
			adultRow = moveColumns16to18After5(adultRow)
			err = adultsFileWriter.Write(adultRow)
			if err != nil {
				fmt.Println("Error writing additional row:", err)
				return
			}
		}

		// Write student fields
		studentIndex++
		studentRow := []string{}
		studentRow = append(studentRow, strconv.Itoa(studentIndex))
		studentRow = append(studentRow, strconv.Itoa(householdIndex))
		studentRow = append(studentRow, row[2])       // full name
		studentRow = append(studentRow, row[3:14]...) // interests
		err = studentsFileWriter.Write(studentRow)
		if err != nil {
			fmt.Println("Error writing row:", err)
			return
		}

		if row[14] == "Yes" {
			studentIndex++
			studentRow := []string{}
			studentRow = append(studentRow, strconv.Itoa(studentIndex))
			studentRow = append(studentRow, strconv.Itoa(householdIndex))
			studentRow = append(studentRow, row[15])       // full name
			studentRow = append(studentRow, row[16:27]...) // interests
			err = studentsFileWriter.Write(studentRow)
			if err != nil {
				fmt.Println("Error writing row:", err)
				return
			}
		}

		if row[27] == "Yes" {
			studentIndex++
			studentRow := []string{}
			studentRow = append(studentRow, strconv.Itoa(studentIndex))
			studentRow = append(studentRow, strconv.Itoa(householdIndex))
			studentRow = append(studentRow, row[28])       // full name
			studentRow = append(studentRow, row[29:40]...) // interests
			err = studentsFileWriter.Write(studentRow)
			if err != nil {
				fmt.Println("Error writing row:", err)
				return
			}
		}

		if row[40] == "Yes" {
			studentIndex++
			studentRow := []string{}
			studentRow = append(studentRow, strconv.Itoa(studentIndex))
			studentRow = append(studentRow, strconv.Itoa(householdIndex))
			studentRow = append(studentRow, row[41])       // full name
			studentRow = append(studentRow, row[42:53]...) // interests
			err = studentsFileWriter.Write(studentRow)
			if err != nil {
				fmt.Println("Error writing row:", err)
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
	cutColumns := row[cutFrom : cutTo+1]

	// Keep the part of the slice before the cut elements
	// And also the part after the cut elements
	sliceBeforeInsert := row[:insertInto]
	sliceAfterInsert := row[insertInto:cutFrom]
	sliceAfterCut := row[cutTo+1:]

	result := make([]string, 0, len(row))
	result = append(result, sliceBeforeInsert...)
	result = append(result, cutColumns...)
	result = append(result, sliceAfterInsert...)
	result = append(result, sliceAfterCut...)

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
