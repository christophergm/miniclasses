package main

import (
	"fmt"
	"log"
)

func main() {
	catalog, err := readClassCatalog("../files/class_catalog.csv")
	if err != nil {
		log.Fatalf("Error reading class catalog: %v", err)
	}

	adults, err := readAdultAssignments("../files/adult_class_assignments.csv")
	if err != nil {
		log.Fatalf("Error reading adult assignments: %v", err)
	}

	students, err := readFinalAssignments("../output/final_assignments.csv")
	if err != nil {
		log.Fatalf("Error reading final assignments: %v", err)
	}

	classData := joinData(catalog, adults, students)

	err = generateMarkdown(classData, "../output/class_list.md")
	if err != nil {
		log.Fatalf("Error generating class list: %v", err)
	}
	fmt.Println("Class list generated successfully.")

	err = generateMarkdownByTeacher(classData, "../output/teacher_list.md")
	if err != nil {
		log.Fatalf("Error generating class list: %v", err)
	}
	fmt.Println("Teacher list generated successfully.")

}
