package main

import (
	"os"
	"sort"
	"strings"
	"text/template"
)

// StudentInfo includes student details along with their class name and location
type StudentInfo struct {
	FinalAssignment
	ClassName         string
	ClassLocation     string
	ClassMeetLocation string
}

// TeacherClassGroup represents a class and its associated students for a teacher
type TeacherClassGroup struct {
	ClassName         string
	ClassLocation     string
	ClassMeetLocation string
	Students          []StudentInfo
}

// Generate a markdown file grouping students by their teacher, and within each teacher, group students by class
func generateMarkdownByTeacher(data map[string]ClassData, outputFile string) error {
	tmpl, err := template.ParseFiles("teacher_list_template.md")
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a map to group classes by teacher, including class details and students
	teacherMap := make(map[string]map[string]*TeacherClassGroup)

	for _, class := range data {
		for _, student := range class.Students {
			teacher := student.StudentTeacher
			if teacherMap[teacher] == nil {
				teacherMap[teacher] = make(map[string]*TeacherClassGroup)
			}

			if _, exists := teacherMap[teacher][class.Catalog.ID]; !exists {
				teacherMap[teacher][class.Catalog.ID] = &TeacherClassGroup{
					ClassName:         class.Catalog.Name,
					ClassLocation:     class.Catalog.Location,
					ClassMeetLocation: class.Catalog.MeetLocation,
					Students:          []StudentInfo{},
				}
			}

			studentInfo := StudentInfo{
				FinalAssignment:   student,
				ClassName:         class.Catalog.Name,
				ClassLocation:     class.Catalog.Location,
				ClassMeetLocation: class.Catalog.MeetLocation,
			}
			teacherMap[teacher][class.Catalog.ID].Students = append(teacherMap[teacher][class.Catalog.ID].Students, studentInfo)
		}
	}

	// Extract and sort teacher names
	teacherNames := make([]string, 0, len(teacherMap))
	for teacher := range teacherMap {
		teacherNames = append(teacherNames, teacher)
	}
	sort.Strings(teacherNames)

	// Render the teachers and their classes with students
	for _, teacher := range teacherNames {
		classMap := teacherMap[teacher]

		// Sort the classes by class ID
		classIDs := make([]string, 0, len(classMap))
		for classID := range classMap {
			classIDs = append(classIDs, classID)
		}
		sort.Strings(classIDs)

		classes := make([]TeacherClassGroup, 0, len(classIDs))
		for _, classID := range classIDs {
			classGroup := classMap[classID]
			// Sort students by name within each class
			sort.Slice(classGroup.Students, func(i, j int) bool {
				return strings.ToLower(classGroup.Students[i].StudentFullName) < strings.ToLower(classGroup.Students[j].StudentFullName)
			})
			classes = append(classes, *classGroup)
		}

		// Render the template for each teacher and their classes
		err := tmpl.Execute(f, struct {
			Teacher string
			Classes []TeacherClassGroup
		}{
			Teacher: teacher,
			Classes: classes,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
