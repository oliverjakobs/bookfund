package main

import (
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
	Amount    float32
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
	stmt, err := db.Prepare("INSERT INTO transactions (type,amount,reason,timestamp) values (?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.Type, t.Amount, t.Reason, t.Timestamp.Unix())

	return err
}

func GetTransactions(after time.Time) ([]Transaction, error) {
	return QueryTransactions("SELECT * FROM transactions WHERE timestamp >= ? ORDER BY timestamp DESC", after.Unix())
}

func GetTransactionsAll() ([]Transaction, error) {
	return QueryTransactions("SELECT * FROM transactions ORDER BY timestamp DESC")
}

func GetBalance() (float32, error) {
	stmt, _ := db.Prepare("SELECT SUM(amount) FROM transactions WHERE type=?")

	var depositAmount, withdrawalAmount float32

	var err error
	err = stmt.QueryRow(DEPOSIT).Scan(&depositAmount)
	err = stmt.QueryRow(WITHDRAWAL).Scan(&withdrawalAmount)

	return depositAmount - withdrawalAmount, err
}
