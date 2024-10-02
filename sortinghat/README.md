# Sorting
This is a simple script for sorting students into mini classes.

## Usage
To generate a sorting, run sort_students.

```shell
$ python sort_students.py
```

## Approach
Given a list of student preferences, we'll sort the students by the pickiest ones first and then
iterate over each student and try to give them a class in one of their top picks.  If we run out 
space, there's a "Fallback" class that will be the catch all for anyone who can't be assigned.