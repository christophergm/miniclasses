package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type Student struct {
	FirstName string
	LastName  string
	Grade     string
	Teacher   string
	Stream    string
	Interests []string
}

func main() {
	// Open the first CSV file
	file1, err := os.Open("student_list.csv")
	if err != nil {
		fmt.Println("Error opening first file:", err)
		return
	}
	defer file1.Close()

	// Open the second CSV file
	file2, err := os.Open("student_preferences.csv")
	if err != nil {
		fmt.Println("Error opening second file:", err)
		return
	}
	defer file2.Close()

	// Read the first CSV file
	reader1 := csv.NewReader(file1)
	records1, err := reader1.ReadAll()
	if err != nil {
		fmt.Println("Error reading first file:", err)
		return
	}

	// Read the second CSV file
	reader2 := csv.NewReader(file2)
	records2, err := reader2.ReadAll()
	if err != nil {
		fmt.Println("Error reading second file:", err)
		return
	}

	// Create a map for the second file where full_name is the key
	studentInterests := make(map[string][]string)
	matched := make(map[string]bool) // To track which records in file2 got matched

	for _, record := range records2[1:] { // Skipping header row
		fullName := strings.TrimSpace(record[2]) // "full_name"
		studentInterests[fullName] = record
		matched[fullName] = false
	}

	// Create a new file for the merged output
	timestamp := time.Now().Format("2006-01-02-1504")
	filename := fmt.Sprintf("student_list_preferences-%s.csv", timestamp)
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write the header row to the new CSV
	header := append(records1[0], records2[0]...) // Combine headers
	writer.Write(header)

	// Process each row of the first file
	for _, record1 := range records1[1:] {
		firstName := strings.TrimSpace(record1[0]) // "first_name"
		lastName := strings.TrimSpace(record1[1])  // "last_name"
		fullName := firstName + " " + lastName

		// Check if the full name exists in the second file's map
		if interests, ok := studentInterests[fullName]; ok {
			// Combine the row with the interests
			combinedRecord := append(record1, interests...)
			writer.Write(combinedRecord)
			// Mark this full_name as matched
			matched[fullName] = true
		} else {
			// Write without interests if no match found
			writer.Write(record1)
		}
	}

	fmt.Println("Merged CSV file created successfully.")

	// Create a file for unmatched entries
	filename = fmt.Sprintf("unmatched-%s.csv", timestamp)
	unmatchedFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating unmatched file:", err)
		return
	}
	defer unmatchedFile.Close()

	// Write to the unmatched CSV file
	unmatchedWriter := csv.NewWriter(unmatchedFile)
	defer unmatchedWriter.Flush()

	// Write the header row for the unmatched CSV (same as file2)
	unmatchedWriter.Write(records2[0])

	// Write all unmatched rows from file2 to the unmatched.csv
	for _, record := range records2[1:] {
		fullName := strings.TrimSpace(record[2])
		if !matched[fullName] {
			// Write the unmatched row to the unmatched.csv
			unmatchedWriter.Write(record)
		}
	}

	fmt.Println("Unmatched rows written to unmatched.csv.")
}
