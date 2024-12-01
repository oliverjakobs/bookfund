sqlite3 db/bookfund-testing.db <<EOF
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY  NOT NULL,
    amount REAL             NOT NULL,
    reason TEXT             NOT NULL,
    timestamp INTEGER       NOT NULL
);
EOF