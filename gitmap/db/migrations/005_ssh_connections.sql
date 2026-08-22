CREATE TABLE IF NOT EXISTS SSHConnection (
    Alias TEXT PRIMARY KEY,
    IPAddress TEXT NOT NULL,
    Username TEXT NOT NULL,
    EncryptedPassword TEXT,
    KeyPath TEXT,
    OS TEXT NOT NULL,
    CreatedAt DATETIME NOT NULL
);
