package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func handleTemplateError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf("index: couldn't parse template: %v", err), http.StatusInternalServerError)
		log.Printf("ERROR: %s\n", err.Error())
	}
}

func FormatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f €", math.Abs(amount))
}

var baseTmpl = template.New("base").Funcs(template.FuncMap{
	"formatCurrency": FormatCurrency,
})

func renderTemplate(w http.ResponseWriter, filename string, data any) {
	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/base.html", filename))

	handleTemplateError(w, tmpl.ExecuteTemplate(w, "base", data))
}

func index(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	timeCutoff := time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, now.Location())
	transactions, err := GetTransactions(timeCutoff)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, _ := GetBalance()

	renderTemplate(w, "./templates/index.html", struct {
		Transactions []Transaction
		Balance      float64
	}{
		Transactions: transactions,
		Balance:      balance,
	})
}

func modal(w http.ResponseWriter, r *http.Request) {
	modalType := r.PathValue("type")

	var templates = template.Must(template.ParseFiles("./templates/modal.html"))
	handleTemplateError(w, templates.ExecuteTemplate(w, "modal", modalType))
}

func transaction(w http.ResponseWriter, r *http.Request) {
	transactionType := TransactionType(r.PathValue("type"))

	err := r.ParseForm()

	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
	}

	var amount float64
	amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)

	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("failed to parse float: %v", err), http.StatusInternalServerError)
	}

	if WITHDRAWAL == transactionType {
		amount *= -1
	}

	err = CreateTransaction(Transaction{
		Type:      transactionType,
		Amount:    amount,
		Reason:    r.FormValue("reason"),
		Timestamp: time.Now(),
	})

	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("failed to create transaction: %v", err), http.StatusInternalServerError)
	}

	w.Header().Add("HX-Redirect", "/")
	w.WriteHeader(http.StatusCreated)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("failed to parse id: %v", err), http.StatusInternalServerError)
	}

	err = DeleteTransaction(id)

	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("failed to delete: %v", err), http.StatusInternalServerError)
	}

	w.Header().Add("HX-Redirect", "/review")
}

func review(w http.ResponseWriter, r *http.Request) {
	transactions, err := GetTransactionsAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/base.html", "./templates/review_table.html", "./templates/review_entries.html"))

	handleTemplateError(w, tmpl.ExecuteTemplate(w, "base", struct {
		Entries []Transaction
	}{
		Entries: transactions,
	}))
}

func reviewSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	transactions, err := GetTransactionsByReason(search)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/review_entries.html"))

	handleTemplateError(w, tmpl.ExecuteTemplate(w, "entries", transactions))
}
