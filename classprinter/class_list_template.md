## {{.Catalog.Name}}

{{- if .Catalog.MeetLocation }}
**Meet at:** {{.Catalog.MeetLocation}}
{{- end }}
**Location:** {{.Catalog.Location}}
**Grades:** {{.Catalog.GradeMin}} - {{.Catalog.GradeMax}}
**Total students:** {{len .Students}}

### Adults
{{range .Adults}}
- {{.FullName}} ({{.Email}}) {{.Note}}
{{end}}

### Students
{{range .Students}}
1. **{{.StudentFullName}}** - Grade {{.StudentGrade}}, {{.StudentTeacher}} ({{.StudentStream}})
{{end}}
