CREATE TABLE IF NOT EXISTS Issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_ref TEXT NULL
    title TEXT,
    description TEXT,
    active INTEGER DEFAULT 0
);


CREATE TABLE IF NOT EXISTS Logs (
    id INTEGER PRIMARY KEY,
    issue_id INTEGER,
    timestamp TEXT,
    entry TEXT,
    FOREIGN KEY(issue_id) REFERENCES Issues(id)
);

CREATE INDEX idx_issues_external_ref ON issues(external_ref);