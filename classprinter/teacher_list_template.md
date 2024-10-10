## {{.Teacher}}

{{range .Classes}}
**{{.ClassName}}** meet at: {{- if .ClassMeetLocation }}{{ .ClassMeetLocation}}{{- else }}{{ .ClassLocation}}{{- end }}

{{range .Students}}
- {{.StudentFullName}} (Grade {{.StudentGrade}})
{{end}}
{{end}}

