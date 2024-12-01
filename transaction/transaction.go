package transaction

import (
	"errors"
	"time"
)

type Transaction struct {
	ID        int64
	Amount    float64
	Reason    string
	Timestamp time.Time
}

func Query(query string, args ...any) ([]Transaction, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return []Transaction{}, err
	}

	var transactions []Transaction
	for rows.Next() {
		var t Transaction

		var unixTime int64
		err = rows.Scan(&t.ID, &t.Amount, &t.Reason, &unixTime)
		if err != nil {
			break
		}
		t.Timestamp = time.Unix(unixTime, 0)

		transactions = append(transactions, t)
	}
	return transactions, err
}

func Create(t Transaction) error {
	_, err := db.Exec("INSERT INTO transactions (amount,reason,timestamp) values (?,?,?)", t.Amount, t.Reason, t.Timestamp.Unix())
	return err
}

func Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	_, err := db.Exec("DELETE FROM transactions WHERE id=?", id)
	return err
}

func GetAfter(after time.Time) ([]Transaction, error) {
	return Query("SELECT * FROM transactions WHERE timestamp >= ? ORDER BY timestamp DESC", after.Unix())
}

func GetAll() ([]Transaction, error) {
	return Query("SELECT * FROM transactions ORDER BY timestamp DESC")
}

func GetByReason(reason string) ([]Transaction, error) {
	return Query("SELECT * FROM transactions WHERE reason LIKE ?  ORDER BY timestamp DESC", "%"+reason+"%")
}

func GetBalance() (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT SUM(amount) FROM transactions").Scan(&balance)

	return balance, err
}
