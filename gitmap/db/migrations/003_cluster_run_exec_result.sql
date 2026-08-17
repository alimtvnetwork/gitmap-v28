CREATE TABLE IF NOT EXISTS ClusterNode (
	NodeId TEXT PRIMARY KEY,
	Alias TEXT NOT NULL,
	DisplayId INTEGER NOT NULL,
	IPAddress TEXT NOT NULL,
	NodeRole TEXT NOT NULL,
	OS TEXT NOT NULL,
	JoinedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	LastHeartbeat DATETIME,
	Status TEXT NOT NULL,
	PasswordHash TEXT,
	PackageManager TEXT
);

CREATE TABLE IF NOT EXISTS ClusterRun (
	ClusterRunId INTEGER PRIMARY KEY AUTOINCREMENT,
	RunRef TEXT NOT NULL,
	CommandKind INTEGER NOT NULL,
	RawCommand TEXT NOT NULL,
	TargetSelector TEXT NOT NULL,
	ExceptClause TEXT,
	StartedAt DATETIME NOT NULL,
	FinishedAt DATETIME,
	TotalNodes INTEGER,
	SucceededNodes INTEGER,
	FailedNodes INTEGER,
	SkippedNodes INTEGER
);

CREATE TABLE IF NOT EXISTS ClusterExecResult (
	ClusterExecResultId INTEGER PRIMARY KEY AUTOINCREMENT,
	ClusterRunId INTEGER NOT NULL,
	NodeId TEXT NOT NULL,
	SubCommand TEXT NOT NULL,
	CommandText TEXT,
	ResultStatus INTEGER NOT NULL,
	ExitCode INTEGER,
	Stdout TEXT,
	Stderr TEXT,
	StartedAt DATETIME,
	FinishedAt DATETIME,
	DurationMs INTEGER,
	ErrorMessage TEXT,
	FOREIGN KEY (ClusterRunId) REFERENCES ClusterRun(ClusterRunId) ON DELETE CASCADE,
	FOREIGN KEY (NodeId) REFERENCES ClusterNode(NodeId) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS IdxClusterRun_StartedAt ON ClusterRun(StartedAt);
CREATE INDEX IF NOT EXISTS IdxClusterExecResult_ClusterRunId ON ClusterExecResult(ClusterRunId);
CREATE INDEX IF NOT EXISTS IdxClusterExecResult_NodeId ON ClusterExecResult(NodeId);
