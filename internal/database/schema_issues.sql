CREATE TABLE IF NOT EXISTS Issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_ref TEXT UNIQUE NULL,
    title TEXT,
    description TEXT,
    active INTEGER DEFAULT 0,
    progress INTEGER
);


CREATE TABLE IF NOT EXISTS Logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issue_id INTEGER,
    timestamp TEXT,
    entry TEXT,
    FOREIGN KEY(issue_id) REFERENCES Issues(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_external_ref ON Issues(external_ref);

DROP VIEW IF EXISTS Active;

CREATE VIEW Active AS 
SELECT id, external_ref, title
FROM Issues i
WHERE i.Active = 1;