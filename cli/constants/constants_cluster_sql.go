package constants

// SQL statements for Cluster tables.
const (
	SQLCreateClusterNode = `CREATE TABLE IF NOT EXISTS ClusterNode (
		NodeId TEXT PRIMARY KEY,
		Alias TEXT NOT NULL DEFAULT "",
		DisplayId INTEGER NOT NULL DEFAULT 0,
		IPAddress TEXT NOT NULL DEFAULT "",
		NodeRole TEXT NOT NULL DEFAULT "client",
		OS TEXT NOT NULL DEFAULT "windows",
		JoinedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		LastHeartbeat TIMESTAMP,
		Status TEXT NOT NULL DEFAULT "online",
		PasswordHash TEXT,
		PackageManager TEXT
	);`

	SQLCreateClusterRun = `CREATE TABLE IF NOT EXISTS ClusterRun (
		ClusterRunId INTEGER PRIMARY KEY AUTOINCREMENT,
		RunRef TEXT NOT NULL UNIQUE,
		CommandKind INTEGER NOT NULL DEFAULT 0,
		RawCommand TEXT NOT NULL DEFAULT "",
		TargetSelector TEXT NOT NULL DEFAULT "",
		ExceptClause TEXT,
		StartedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FinishedAt TIMESTAMP,
		TotalNodes INTEGER,
		SucceededNodes INTEGER,
		FailedNodes INTEGER,
		SkippedNodes INTEGER
	);`

	SQLCreateClusterExecResult = `CREATE TABLE IF NOT EXISTS ClusterExecResult (
		ClusterExecResultId INTEGER PRIMARY KEY AUTOINCREMENT,
		ClusterRunId INTEGER NOT NULL,
		NodeId TEXT NOT NULL,
		SubCommand TEXT NOT NULL,
		CommandText TEXT NOT NULL,
		ResultStatus INTEGER NOT NULL DEFAULT 0,
		Stdout TEXT,
		Stderr TEXT,
		ExitCode INTEGER,
		ExecutedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		DurationMs INTEGER,
		FOREIGN KEY(ClusterRunId) REFERENCES ClusterRun(ClusterRunId) ON DELETE CASCADE
	);`
)
