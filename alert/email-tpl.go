package alert

import (
	"bytes"
	"text/template"
	"time"
)

func getEmailTemplate(alertList []string) (string, string) {

	templateBody := `The following alerts have been registered:
{{range .alertList}} 
{{.}} {{end}}
	
This email has been automatically generated at {{.time}}`

	templateSubject := `[Prosit notification] Alert`

	bodyData := make(map[string]interface{})
	bodyData["alertList"] = alertList
	bodyData["time"] = time.Now().String()

	tmpl, err := template.New("email").Parse(templateBody)
	if err != nil {
		panic(err)
	}

	body := &bytes.Buffer{}
	err = tmpl.Execute(body, bodyData)
	if err != nil {
		panic(err)
	}

	return templateSubject, string(body.Bytes())

}
