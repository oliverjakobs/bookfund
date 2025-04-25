package main

import (
	"fmt"
	"html/template"
	"math"
)

func FormatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f €", amount)
}

var baseTmpl = template.New("base").Funcs(template.FuncMap{
	"formatCurrency": FormatCurrency,
	"abs":            math.Abs,
})

func CreateTemplate(filenames ...string) *template.Template {
	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles(filenames...))

	return tmpl
}
