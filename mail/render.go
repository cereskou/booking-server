package mail

import (
	"bytes"
	"ditto/booking/utils"
	"text/template"
)

//Render -
func (m *Mail) Render(id string, tmpl string, data map[string]interface{}) (string, error) {
	tmp, err := utils.JSON.Marshal(data)
	if err != nil {
		return tmpl, err
	}
	//val
	var val interface{}
	err = utils.JSON.Unmarshal(tmp, &val)
	if err != nil {
		return tmpl, err
	}

	var tt *template.Template
	tt = template.New(id)
	t, err := tt.Parse(tmpl)
	if err != nil {
		return tmpl, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, val)
	if err != nil {
		return tmpl, err
	}

	return buf.String(), nil
}
