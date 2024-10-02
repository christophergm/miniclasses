from collections import defaultdict
import os
import csv
from typing import Any, Dict, List, Set
from enum import IntEnum
from random import shuffle

join_path = os.path.join

data_dir = join_path(os.path.dirname(__file__), '../files_test')

CURRENT_SESSION = 1

# TODO
# - Manual assignments are not honored yet.
# - We aren't outputing the CSV of assignments yet
# - Paths are hard coded (this can wait)

# List of student preferences
# CSV columns: student_full_name,student_interest_games_puzzles,student_interest_arts_crafts,student_interest_performing_arts,student_interest_cooking,student_interest_athletics,student_interest_building_making,student_interest_gardening,student_interest_science_nature,student_interest_community,student_interest_fabric_arts,student_interest_book_club
student_preferences_path =  join_path(data_dir,'student_preferences.csv')

# List of all students
# CSV columns: first_name,last_name,grade,teacher,stream
student_list_path =  join_path(data_dir,'student_list.csv')

# List of all courses
# CSV columns: id,session,name,interest_area,grade_min,grade_max,student_capacity_max
course_list_path =  join_path(data_dir,'class_catalog.csv')

class Interest(IntEnum): 
  VERY = 0
  MAYBE = 1
  NOPE = 2


class Preference:
  area: str
  level: Interest

  def __init__(self, area: str, level: str):
    self.area = area
    if level == "Very Interested":
      self.level = Interest.VERY
    elif level == "Interested":
      self.level = Interest.MAYBE
    else:
      self.level = Interest.NOPE
  
  def __repr__(self) -> str:
    return f"{self.area} {self.level.name}"


class Student:
  name: str
  grade: int
  _preferences: List[Preference]
  preferences_by_interest: Dict[Interest, List[Preference]]
  
  def __init__(self, row: Dict[str, str]):
    self.name = f"{row['first_name']} {row['last_name']}"
    self.grade = int(row['grade'])

  @property
  def preferences(self):
    return self._preferences

  @preferences.setter
  def preferences(self, prefs: List[Preference]):
    # Shuffle the preferences so that if a student has multiple
    # options we don't always pick them in the same order.
    shuffle(prefs)
    self._preferences = prefs
    self.preferences_by_interest = defaultdict(list)
    for p in prefs:
      self.preferences_by_interest[p.level].append(p)

  def get_preference_counts(self):
    very_count = len(self.preferences_by_interest[Interest.VERY])
    maybe_count = len(self.preferences_by_interest[Interest.MAYBE])
    nope_count = len(self.preferences_by_interest[Interest.NOPE])
    return very_count, maybe_count, nope_count
  
  def get_ordered_preferences(self):
    very = self.preferences_by_interest[Interest.VERY]
    maybe = self.preferences_by_interest[Interest.MAYBE]
    nope = self.preferences_by_interest[Interest.NOPE]
    return very + maybe + nope
  
  def interest_in_course(self, course: "Course"):
    for pref in self.preferences:
      if pref.area == course.area:
        return pref.level
    return Interest.NOPE

  def __repr__(self) -> str:
    return f"{self.name} ({self.grade})"


class Course:
  name: str
  area: str
  grade_min: int
  grade_max: int
  max_capacity: int 
  capacity: int
  students: List[Student]
  
  def __init__(self, row: Dict[str, str]):
    self.name = row['name']
    self.area = row['interest_area']
    self.grade_min = int(row['grade_min'])
    self.grade_max = int(row['grade_max'])
    self.capacity = int(row['student_capacity_max'])
    self.max_capacity = self.capacity
    self.students = []

  def available_to(self, student: Student):
    if self.capacity <= 0:
      return False
    
    if student.grade < self.grade_min or student.grade > self.grade_max:
      return False
    
    return True

  def assign(self, student: Student):
    self.students.append(student)
    self.capacity -= 1

  def __repr__(self) -> str:
    return f"{self.name} ({self.area})"


def interest_area(column_name: str, prefix="student_interest_"):
  """
  Helper to pluck out interest areas from column names
  
  The expected format is "student_interest_computer_typing"
  """
  if column_name.startswith(prefix):
        return column_name.replace(prefix, "")
  
  return None


# Read in the input files and assemble all of our data.

# Start with the student preferences.
preferences_by_student: Dict[str, List[Preference]]  = {}
DEFAULT_PREFERENCES = []
KNOWN_AREAS: Set[str] = set()

with open(student_preferences_path) as csvfile:
  reader = csv.DictReader(csvfile)

  # We populate the DEFAULT_PREFERENCES form the header row.
  #
  # We assume that they were so wildly interested by every area that
  # they simply couldn't decide. This simplifies the algorithm by sorting
  # them to the end of the list. 
  # 
  # Sorry folks, sign up next session!
  for key in reader.fieldnames:
    area = interest_area(key)
    if area:
      KNOWN_AREAS.add(area)
      DEFAULT_PREFERENCES.append(Preference(area, "Very Interested"))

  # Generate a Prefererence for every signup.
  for row in reader:
    preferences = []
    for key, value in row.items():
      area = interest_area(key)
      if area:
        preferences.append(Preference(area, value))
    
    name = row['student_full_name']
    preferences_by_student[name] = preferences

# Next we read in the students and then match them up with their preferences
students: List[Student] = []
with open(student_list_path) as csvfile:
  reader = csv.DictReader(csvfile)
  for row in reader:
    student = Student(row)
    preferences = preferences_by_student.get(student.name)
    if preferences:
      student.preferences = preferences
    else:
      print(f"No preferences for {student.name}")
      student.preferences = DEFAULT_PREFERENCES[:]

    students.append(student)

# Finally, read in all the courses.
DEFAULT_COURSE = Course({
  "name": "Fallback",
  "interest_area": "none",
  "grade_min": 0,
  "grade_max": 999,
  "student_capacity_max": 999
})
courses: List[Course] = [DEFAULT_COURSE]
courses_by_area: Dict[str, List[Course]] = defaultdict(list)
with open(course_list_path) as csvfile:
  reader = csv.DictReader(csvfile)
  for row in reader:
    session = int(row['session'])
    if session == CURRENT_SESSION:
      course = Course(row)
      assert course.area in KNOWN_AREAS, f"Unexpected interest area ({course.area}) for course {course.name}"
      courses.append(course)
      courses_by_area[course.area].append(course)


print(f"Assinging {len(students)} to {len(courses)} courses.")

# Sort the students by the pickiest (fewest VERY interested areas)
students.sort(key=Student.get_preference_counts)

# We'll assign students by iterating over their preferences and
# then checking to see which classes they could be assigned to.
def assign_student(student: Student):
  for pref in student.get_ordered_preferences():
    candidate_courses = courses_by_area[pref.area]
    
    # Shuffle the courses each time so we don't always assign
    # people to the first course listed.
    shuffle(candidate_courses)

    for course in candidate_courses:
      if course.available_to(student):
        course.assign(student)
        return course
      
  print(f"No available class for {student}")
  DEFAULT_COURSE.assign(student)
  return DEFAULT_COURSE

for student in students:
  course = assign_student(student)
  print(f"Assigning {student} to {course}")

# And we're done.  Now we just need to print out the courses!
for course in courses:
  print(f"{course.name} ({len(course.students)}/{course.max_capacity})")
  for student in course.students:
    print(f"- {student} ({student.interest_in_course(course).name})")
