package alert

import (
	"bytes"
	"text/template"
	"time"
)

var templateBody string = `The following alerts have been registered:
{{range .alertList}} 
{{.}} {{end}}

This email has been automatically generated at {{.time}}`

var templateSubject string = `[Prosit notification] Alert: '{{.firstAlert}}'`

func getEmailTemplate(alertList []string) (string, string) {

	emailData := make(map[string]interface{})
	emailData["alertList"] = alertList
	emailData["firstAlert"] = alertList[0]
	emailData["time"] = time.Now().String()

	return processTemplate(templateSubject, emailData), processTemplate(templateBody, emailData)
}

func processTemplate(templateStr string, data interface{}) string {

	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		panic(err)
	}

	body := &bytes.Buffer{}
	err = tmpl.Execute(body, data)
	if err != nil {
		panic(err)
	}

	return string(body.Bytes())
}
