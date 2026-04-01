CREATE TABLE IF NOT EXISTS Issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_ref TEXT NULL,
    title TEXT,
    description TEXT,
    active INTEGER DEFAULT 0
);


CREATE TABLE IF NOT EXISTS Logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issue_id INTEGER,
    timestamp TEXT,
    entry TEXT,
    FOREIGN KEY(issue_id) REFERENCES Issues(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_external_ref ON Issues(external_ref);