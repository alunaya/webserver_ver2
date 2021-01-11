package cws

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Context struct {
	res http.ResponseWriter
	req *http.Request
	templates *template.Template
}

func NewContext(w http.ResponseWriter, r *http.Request, templates *template.Template) *Context{
	return &Context{
		res: w,
		req: r,
		templates: templates,
	}
}

func (c Context)JsonResult(statusCode int, body interface{}) {
	if body != nil {
		jsonResult, err := json.Marshal(body)
		if err != nil {
			http.Error(c.res, err.Error(), http.StatusInternalServerError)
			return
		}

		_, writeError := c.res.Write(jsonResult)

		if writeError != nil {
			http.Error(c.res, err.Error(), http.StatusInternalServerError)
		}

		c.res.WriteHeader(statusCode)
	}
}

func (c Context) StreamResult() {
	panic("not implemented")
}

func (c* Context) Page(fileName string, data interface{}) {
	parseError := c.templates.ExecuteTemplate(c.res, fileName, data)
	if parseError != nil{
		http.Error(c.res, "webpage missing", 404)
	}
}