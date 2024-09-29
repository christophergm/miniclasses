# OPTO Mini Class Planner

# Overview

This contains source files and an algorithm to sort students using preference weighting into classes.

The folder `files_test` contains mock data for development and testing.

# Student Assignment Process

Class assignments are made in the following steps:

### 1. Clean form data

Raw data from Google form submissions is scrubbed to ensure student full names are correct (e.g. typo in name) and complete (e.g. only student first name is entered).

A new header row is added to the Google form data to make more parseable column keys.

### 2. Collect student preferences and class catalog data

The following files are inputs:

* `student_list.csv` 

    - static list of students pulled from the school directory

* `sign_up_form_YYYY-MM-DD.csv`
     - Google form responses submitted by parents
     - each row can have multiple children
     - adult participation preferences are also included but are not used for the class preferences

* `class_catalog.csv`
    - list of class that will be offered for the session
    - each class lists the grade range of students that are eligible to join and the maximum number of students

* `class_assignments_manual.csv`
    - list of students that will be manually assigned to classes by a decision of the organizer and won't be assigned with the sorting algorithm. 
    - this supports cases like where an adult Class Leader and their child want to be together


### 3. Apply assignment algorithm

Algorithm TBD, but will do something like the following.

Randomly assign students to classes with weighting for:

1. students in classes of topics that they have a high interest in
2. prioritizing preferences of students that didn't get a high preference topic in a previous session (not relevant for first session)
2. achieving a mix of Green and Blue streams

and constrain assignments by the following:

1. class can have no more students than their maximum capacity
2. students can only be assigned to classes within their grade range (unless they are manually assigned)

Students that did not have a completed preference form should still be included in a class assignment.

### 4. Save final assignments

* `class_assignments_final.csv`

    - result of the assignment algorithm is a list of each student and their class assignment
