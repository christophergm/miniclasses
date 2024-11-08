package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type ClassCatalog struct {
	ID              string
	Session         int
	Name            string
	InterestArea    string
	GradeMin        int
	GradeMax        int
	StudentCapacity int
	Location        string
	MeetLocation    string
}

type AdultClassAssignment struct {
	ClassID  string
	FullName string
	Email    string
	Note     string
}

type FinalAssignment struct {
	ClassName       string
	ClassSession    int
	ClassID         string
	StudentFullName string
	StudentGrade    int
	StudentTeacher  string
	StudentStream   string
	StudentInterest string
}

func readClassCatalog(file string) ([]ClassCatalog, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var catalog []ClassCatalog
	for _, record := range records[1:] { // Skipping header row
		session, _ := strconv.Atoi(record[1])
		gradeMin, _ := strconv.Atoi(record[4])
		gradeMax, _ := strconv.Atoi(record[5])
		capacity, _ := strconv.Atoi(record[6])

		catalog = append(catalog, ClassCatalog{
			ID:              record[0],
			Session:         session,
			Name:            record[2],
			InterestArea:    record[3],
			GradeMin:        gradeMin,
			GradeMax:        gradeMax,
			StudentCapacity: capacity,
			Location:        record[7],
			MeetLocation:    record[8],
		})
	}

	return catalog, nil
}

func readAdultAssignments(file string) ([]AdultClassAssignment, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var assignments []AdultClassAssignment
	for _, record := range records[1:] {
		assignments = append(assignments, AdultClassAssignment{
			ClassID:  record[0],
			FullName: record[1],
			Email:    record[2],
			Note:     record[3],
		})
	}

	return assignments, nil
}

func readFinalAssignments(file string) ([]FinalAssignment, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var assignments []FinalAssignment
	for _, record := range records[1:] {
		session, _ := strconv.Atoi(record[1])
		grade, _ := strconv.Atoi(record[4])

		assignments = append(assignments, FinalAssignment{
			ClassName:       record[0],
			ClassSession:    session,
			ClassID:         record[2],
			StudentFullName: record[3],
			StudentGrade:    grade,
			StudentTeacher:  record[5],
			StudentStream:   record[6],
			StudentInterest: record[7],
		})
	}

	return assignments, nil
}
