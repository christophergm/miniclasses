package main

import (
	"os"
	"sort"
	"strings"
	"text/template"
)

type ClassData struct {
	Catalog  ClassCatalog
	Adults   []AdultClassAssignment
	Students []FinalAssignment
}

func joinData(catalog []ClassCatalog, adults []AdultClassAssignment, students []FinalAssignment) map[string]ClassData {
	classMap := make(map[string]ClassData)

	// Build initial map from catalog
	for _, class := range catalog {
		classMap[class.ID] = ClassData{Catalog: class}
	}

	// Add adults to the map
	for _, adult := range adults {
		if class, exists := classMap[adult.ClassID]; exists {
			class.Adults = append(class.Adults, adult)
			classMap[adult.ClassID] = class
		}
	}

	// Add students to the map
	for _, student := range students {
		if class, exists := classMap[student.ClassID]; exists {
			class.Students = append(class.Students, student)
			classMap[student.ClassID] = class
		}
	}

	return classMap
}

// Sort adults alphabetically by full name
func sortAdultsByName(adults []AdultClassAssignment) {
	sort.Slice(adults, func(i, j int) bool {
		return strings.ToLower(adults[i].FullName) < strings.ToLower(adults[j].FullName)
	})
}

// Sort students by grade (ascending) and then by first name (alphabetically)
func sortStudentsByGradeAndName(students []FinalAssignment) {
	sort.Slice(students, func(i, j int) bool {
		if students[i].StudentGrade == students[j].StudentGrade {
			return strings.ToLower(strings.Split(students[i].StudentFullName, " ")[0]) <
				strings.ToLower(strings.Split(students[j].StudentFullName, " ")[0])
		}
		return students[i].StudentGrade < students[j].StudentGrade
	})
}
func generateMarkdown(data map[string]ClassData, outputFile string) error {
	tmpl, err := template.ParseFiles("class_list_template.md")
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Extract and sort class IDs
	classIDs := make([]string, 0, len(data))
	for classID := range data {
		classIDs = append(classIDs, classID)
	}
	sort.Strings(classIDs) // Sort classIDs alphabetically

	// Render the classes in sorted order
	for _, classID := range classIDs {
		class := data[classID]

		// Sort adults by name
		sortAdultsByName(class.Adults)

		// Sort students by grade, then first name
		sortStudentsByGradeAndName(class.Students)

		err := tmpl.Execute(f, class)
		if err != nil {
			return err
		}
	}

	return nil
}
