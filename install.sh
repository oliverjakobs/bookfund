sqlite3 db/bookfund-develop.db <<EOF
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY  NOT NULL,
    type TEXT               NOT NULL,
    amount REAL             NOT NULL,
    reason TEXT             NOT NULL,
    timestamp INTEGER       NOT NULL
);

INSERT INTO transactions (type,amount,reason,timestamp) values ('deposit',30.0,'Oktober',1727733600);
INSERT INTO transactions (type,amount,reason,timestamp) values ('withdrawal',8.0,'Windhaven',1729029600);
INSERT INTO transactions (type,amount,reason,timestamp) values ('withdrawal',20.0,'Die Großen Religionen der Welt',1729893600);
INSERT INTO transactions (type,amount,reason,timestamp) values ('deposit',30.0,'November',1730415600);
EOF