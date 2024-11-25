package main

import (
	t "bookfund/transaction"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func handleError(w http.ResponseWriter, msg string, err error) {
	http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusInternalServerError)
	log.Printf("ERROR: %s\n", err.Error())
}

func FormatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f €", math.Abs(amount))
}

var baseTmpl = template.New("base").Funcs(template.FuncMap{
	"formatCurrency": FormatCurrency,
})

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	err := tmpl.ExecuteTemplate(w, name, data)

	if err != nil {
		handleError(w, "failed to render template", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	timeCutoff := time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, now.Location())
	transactions, err := t.GetAfter(timeCutoff)

	if err != nil {
		handleError(w, "failed to query transactions", err)
		return
	}

	var balance float64
	balance, err = t.GetBalance()
	if err != nil {
		handleError(w, "failed to query balance", err)
		return
	}

	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/base.html", "./templates/index.html"))

	renderTemplate(w, tmpl, "base", struct {
		Transactions []t.Transaction
		Balance      float64
	}{
		Transactions: transactions,
		Balance:      balance,
	})
}

func modal(w http.ResponseWriter, r *http.Request) {
	modalType := r.PathValue("type")

	tmpl := template.Must(template.ParseFiles("./templates/modal.html"))

	renderTemplate(w, tmpl, "modal", modalType)
}

func parseTransactionForm(r *http.Request, transactionType t.TransactionType) (t.Transaction, error) {
	err := r.ParseForm()

	if err != nil {
		return t.Transaction{}, err
	}

	var amount float64
	amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)

	if err != nil {
		return t.Transaction{}, err
	}

	if t.WITHDRAWAL == transactionType {
		amount *= -1
	}

	return t.Transaction{
		Type:      transactionType,
		Amount:    amount,
		Reason:    r.FormValue("reason"),
		Timestamp: time.Now(),
	}, nil
}

func post(w http.ResponseWriter, r *http.Request) {
	transactionType := t.TransactionType(r.PathValue("type"))

	transaction, err := parseTransactionForm(r, transactionType)

	if err != nil {
		handleError(w, "failed to parse form", err)
		return
	}

	err = t.Create(transaction)

	if err != nil {
		handleError(w, "failed to create transaction", err)
		return
	}

	w.Header().Add("HX-Redirect", "/")
	w.WriteHeader(http.StatusCreated)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		handleError(w, "failed to parse id", err)
		return
	}

	err = t.Delete(id)

	if err != nil {
		handleError(w, "failed to delete transaction", err)
		return
	}

	w.Header().Add("HX-Redirect", "/review")
}

func review(w http.ResponseWriter, r *http.Request) {
	transactions, err := t.GetAll()

	if err != nil {
		handleError(w, "failed to query transactions", err)
		return
	}

	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/base.html", "./templates/review_table.html", "./templates/review_entries.html"))

	renderTemplate(w, tmpl, "base", struct {
		Entries []t.Transaction
	}{
		Entries: transactions,
	})
}

func reviewSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	transactions, err := t.GetByReason(search)

	if err != nil {
		handleError(w, "failed to query transactions", err)
		return
	}

	tmpl := template.Must(baseTmpl.Clone())
	template.Must(tmpl.ParseFiles("./templates/review_entries.html"))

	renderTemplate(w, tmpl, "entries", transactions)
}
