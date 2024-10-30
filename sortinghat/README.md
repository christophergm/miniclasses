# Sorting
This is a simple script for sorting students into mini classes.

## Usage
Step 1: Setup input CSV files (see source code for required columns)


* `student_preferences.csv` - List of student preferences
* `student_list.csv` - List of all students
* `class_catalog.csv` - List of all courses
* `class_assignments_manual.csv` - Manual class assignments
* `class_exclusions_manual.csv` - Manual class exclusions
* `skip_assignments_manual.csv` - Exclude some students from the student list (e.g. wrong grade)

Step 2: Generate a sorting by running `sort_students`.

```shell
$ python sort_students.py
```

## Approach
Given a list of student preferences, we'll sort the students by the pickiest ones first and then
iterate over each student and try to give them a class in one of their top picks.  If we run out 
space, there's a "Fallback" class that will be the catch all for anyone who can't be assigned.