package main

import (
	"errors"
	"time"
)

type TransactionType string

const (
	DEPOSIT    = "deposit"
	WITHDRAWAL = "withdrawal"
)

type Transaction struct {
	ID        int64
	Type      TransactionType
	Amount    float64
	Reason    string
	Timestamp time.Time
}

func QueryTransactions(query string, args ...any) ([]Transaction, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return []Transaction{}, err
	}

	var transactions []Transaction
	for rows.Next() {
		var t Transaction

		var unixTime int64
		err = rows.Scan(&t.ID, &t.Type, &t.Amount, &t.Reason, &unixTime)
		if err != nil {
			break
		}
		t.Timestamp = time.Unix(unixTime, 0)

		transactions = append(transactions, t)
	}
	return transactions, err
}

func CreateTransaction(t Transaction) error {
	_, err := db.Exec("INSERT INTO transactions (type,amount,reason,timestamp) values (?,?,?,?)", t.Type, t.Amount, t.Reason, t.Timestamp.Unix())
	return err
}

func DeleteTransaction(id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	_, err := db.Exec("DELETE FROM transactions WHERE id=?", id)
	return err
}

func GetTransactions(after time.Time) ([]Transaction, error) {
	return QueryTransactions("SELECT * FROM transactions WHERE timestamp >= ? ORDER BY timestamp DESC", after.Unix())
}

func GetTransactionsAll() ([]Transaction, error) {
	return QueryTransactions("SELECT * FROM transactions ORDER BY timestamp DESC")
}

func GetTransactionsByReason(reason string) ([]Transaction, error) {
	return QueryTransactions("SELECT * FROM transactions WHERE reason LIKE ?  ORDER BY timestamp DESC", "%"+reason+"%")
}

func GetBalance() (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT SUM(amount) FROM transactions").Scan(&balance)

	return balance, err
}
