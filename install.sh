sqlite3 db/bookfund-testing.db <<EOF
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY  NOT NULL,
    type TEXT               NOT NULL,
    amount REAL             NOT NULL,
    reason TEXT             NOT NULL,
    timestamp INTEGER       NOT NULL
);
EOF