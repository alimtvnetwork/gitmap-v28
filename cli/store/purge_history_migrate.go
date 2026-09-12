package store

// EnsurePurgeHistoryTable ensures that the PurgeHistoryLog table exists and its columns are migrated.
func (db *DB) EnsurePurgeHistoryTable() error {
	if _, err := ExecWrapper(db.conn, sqlCreatePurgeHistory).Destruct(); err != nil {
		return err
	}

	return db.migratePurgeHistoryColumns()
}

// migratePurgeHistoryColumns renames legacy ID -> PurgeHistoryLogId and Restored -> IsRestored,
// and adds Notes and Comments columns if they do not exist.
func (db *DB) migratePurgeHistoryColumns() error {
	if !db.tableExists("PurgeHistoryLog") {
		return nil
	}

	if err := db.migratePurgeHistoryIdColumn(); err != nil {
		return err
	}

	if err := db.migratePurgeHistoryRestoredColumn(); err != nil {
		return err
	}

	return db.migratePurgeHistoryContextColumns()
}

func (db *DB) execAlterIf(isNeeded bool, sql string) error {
	if !isNeeded {
		return nil
	}

	_, err := ExecWrapper(db.conn, sql).Destruct()
	if err != nil && !isBenignAlterError(err) {
		return err
	}

	return nil
}

func (db *DB) migratePurgeHistoryIdColumn() error {
	if db.columnExists("PurgeHistoryLog", "PurgeHistoryLogId") {
		return nil
	}

	oldCol := ""
	if db.columnExists("PurgeHistoryLog", "ID") {
		oldCol = "ID"
	} else if db.columnExists("PurgeHistoryLog", "Id") {
		oldCol = "Id"
	}

	sql := "ALTER TABLE PurgeHistoryLog RENAME COLUMN " + oldCol + " TO PurgeHistoryLogId"

	return db.execAlterIf(oldCol != "", sql)
}

func (db *DB) migratePurgeHistoryRestoredColumn() error {
	if db.columnExists("PurgeHistoryLog", "IsRestored") {
		return nil
	}

	hasOld := db.columnExists("PurgeHistoryLog", "Restored")
	sql := "ALTER TABLE PurgeHistoryLog RENAME COLUMN Restored TO IsRestored"

	return db.execAlterIf(hasOld, sql)
}

func (db *DB) migratePurgeHistoryContextColumns() error {
	isNotesMissing := !db.columnExists("PurgeHistoryLog", "Notes")
	if err := db.execAlterIf(isNotesMissing, "ALTER TABLE PurgeHistoryLog ADD COLUMN Notes TEXT NULL"); err != nil {
		return err
	}

	isCommentsMissing := !db.columnExists("PurgeHistoryLog", "Comments")

	return db.execAlterIf(isCommentsMissing, "ALTER TABLE PurgeHistoryLog ADD COLUMN Comments TEXT NULL")
}
