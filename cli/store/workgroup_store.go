package store

import (
	"database/sql"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func scanWorkDirectoryRow(scanner interface{ Scan(...any) error }) (*model.WorkDirectory, error) {
	var item model.WorkDirectory
	var isDef, hasBinding int

	err := scanner.Scan(&item.WorkDirectoryId, &item.AbsolutePath, &item.Label, &isDef, &hasBinding, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scanWorkDirectoryRow")
	}

	item.IsDefault = (isDef == 1)
	item.HasActiveBinding = (hasBinding == 1)

	return &item, nil
}

func scanFolderWorkDirRow(scanner interface{ Scan(...any) error }) (*model.FolderWorkDir, error) {
	var item model.FolderWorkDir
	var isRec, isAct int

	err := scanner.Scan(&item.FolderWorkDirId, &item.FolderPath, &item.WorkDirectoryId, &isRec, &item.Priority, &isAct, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scanFolderWorkDirRow")
	}

	item.IsRecursive = (isRec == 1)
	item.IsActive = (isAct == 1)

	return &item, nil
}

func scanFolderGroupRow(scanner interface{ Scan(...any) error }) (*model.FolderGroup, error) {
	var item model.FolderGroup
	var isDef, isAct int

	err := scanner.Scan(&item.FolderGroupId, &item.GroupName, &item.FolderPath, &item.Description, &isDef, &isAct, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scanFolderGroupRow")
	}

	item.IsDefault = (isDef == 1)
	item.IsActive = (isAct == 1)

	return &item, nil
}

func scanSshKeyAssociationRow(scanner interface{ Scan(...any) error }) (*model.SshKeyAssociation, error) {
	var item model.SshKeyAssociation
	var isAct, isAuto int

	err := scanner.Scan(&item.SshKeyAssociationId, &item.SshKeyId, &item.AssociationType, &item.TargetIdentifier, &item.FolderPath, &item.GroupName, &item.WorkDirectoryId, &item.RepoId, &item.GitAuthorName, &item.GitAuthorEmail, &isAct, &isAuto, &item.Priority, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scanSshKeyAssociationRow")
	}

	item.IsActive = (isAct == 1)
	item.IsAutoInjected = (isAuto == 1)

	return &item, nil
}

// EnsureWorkDirectory registers or updates a WorkDirectory in the database.
func (db *DB) EnsureWorkDirectory(absPath, label string, isDefault bool) (*model.WorkDirectory, error) {
	isDefInt := 0
	if isDefault {
		isDefInt = 1
	}

	_, errExec := db.conn.Exec(constants.SQLUpsertWorkDirectory, absPath, label, isDefInt, 1)
	if errExec != nil {
		return nil, apperror.WrapSimple(errExec, "EnsureWorkDirectory.Exec")
	}

	return db.GetWorkDirectoryByPath(absPath)
}

// GetWorkDirectoryByPath returns the WorkDirectory for the given absolute path.
func (db *DB) GetWorkDirectoryByPath(absPath string) (*model.WorkDirectory, error) {
	row := db.conn.QueryRow(constants.SQLSelectWorkDirectoryByPath, absPath)
	item, errScan := scanWorkDirectoryRow(row)
	if errScan != nil {
		return nil, errScan
	}

	return item, nil
}

// GetWorkDirectoryById returns the WorkDirectory for the given identifier.
func (db *DB) GetWorkDirectoryById(id int64) (*model.WorkDirectory, error) {
	row := db.conn.QueryRow(constants.SQLSelectWorkDirectoryById, id)
	item, errScan := scanWorkDirectoryRow(row)
	if errScan != nil {
		return nil, errScan
	}

	return item, nil
}

// ListWorkDirectories returns all registered work directories.
func (db *DB) ListWorkDirectories() ([]model.WorkDirectory, error) {
	rows, errQuery := db.conn.Query(constants.SQLSelectAllWorkDirectories)
	if errQuery != nil {
		return nil, apperror.WrapSimple(errQuery, "ListWorkDirectories.Query")
	}

	defer rows.Close()

	return collectWorkDirectories(rows)
}

func collectWorkDirectories(rows *sql.Rows) ([]model.WorkDirectory, error) {
	var list []model.WorkDirectory
	for rows.Next() {
		item, errScan := scanWorkDirectoryRow(rows)
		if errScan != nil {
			return nil, errScan
		}

		list = append(list, *item)
	}

	return list, nil
}

// SetDefaultWorkDirectory marks the target WorkDirectory as default.
func (db *DB) SetDefaultWorkDirectory(idOrPath string) error {
	_, err := db.conn.Exec(constants.SQLSetDefaultWorkDirectory, idOrPath, idOrPath)
	if err != nil {
		return apperror.WrapSimple(err, "SetDefaultWorkDirectory")
	}

	return nil
}

// DeleteWorkDirectory removes a WorkDirectory by ID or absolute path.
func (db *DB) DeleteWorkDirectory(idOrPath string) error {
	_, err := db.conn.Exec(constants.SQLDeleteWorkDirectory, idOrPath, idOrPath)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteWorkDirectory")
	}

	return nil
}

func (db *DB) getFolderWorkDirByPath(folderPath string) (*model.FolderWorkDir, error) {
	const sqlFind = "SELECT FolderWorkDirId, FolderPath, WorkDirectoryId, IsRecursive, Priority, IsActive, CreatedAt, UpdatedAt FROM FolderWorkDir WHERE FolderPath = ? COLLATE NOCASE"
	row := db.conn.QueryRow(sqlFind, folderPath)
	item, err := scanFolderWorkDirRow(row)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// BindFolderWorkDir associates a folder path with a WorkDirectory.
func (db *DB) BindFolderWorkDir(folderPath string, workDirId int64, isRecursive bool, priority int) (*model.FolderWorkDir, error) {
	isRecInt := 0
	if isRecursive {
		isRecInt = 1
	}

	_, err := db.conn.Exec(constants.SQLUpsertFolderWorkDir, folderPath, workDirId, isRecInt, priority, 1)
	if err != nil {
		return nil, apperror.WrapSimple(err, "BindFolderWorkDir.Exec")
	}

	return db.getFolderWorkDirByPath(folderPath)
}

// UnbindFolderWorkDir removes the binding for the specified folder path.
func (db *DB) UnbindFolderWorkDir(folderPath string) error {
	_, err := db.conn.Exec(constants.SQLDeleteFolderWorkDir, folderPath)
	if err != nil {
		return apperror.WrapSimple(err, "UnbindFolderWorkDir")
	}

	return nil
}

// ListFolderWorkDirs returns all folder-to-workdir bindings.
func (db *DB) ListFolderWorkDirs() ([]model.FolderWorkDir, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllFolderWorkDirs)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListFolderWorkDirs.Query")
	}

	defer rows.Close()

	return collectFolderWorkDirs(rows)
}

func collectFolderWorkDirs(rows *sql.Rows) ([]model.FolderWorkDir, error) {
	var list []model.FolderWorkDir
	for rows.Next() {
		item, err := scanFolderWorkDirRow(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *item)
	}

	return list, nil
}

type joinedWorkDirRow struct {
	fId, wId, isRec, prio, isAct, wId2, isDef, hasBind int64
	fPath, wAbs, wLabel, cAt, uAt                      string
}

func (r *joinedWorkDirRow) toWorkDirectory() *model.WorkDirectory {
	return &model.WorkDirectory{
		WorkDirectoryId:  r.wId2,
		AbsolutePath:     r.wAbs,
		Label:            r.wLabel,
		IsDefault:        (r.isDef == 1),
		HasActiveBinding: (r.hasBind == 1),
		CreatedAt:        r.cAt,
		UpdatedAt:        r.uAt,
	}
}

func queryNearestFolderWorkDir(db *DB, targetPath string) (*model.WorkDirectory, error) {
	var r joinedWorkDirRow
	row := db.conn.QueryRow(constants.SQLResolveWorkDirectoryForPath, targetPath, targetPath, targetPath)
	err := row.Scan(&r.fId, &r.fPath, &r.wId, &r.isRec, &r.prio, &r.isAct, &r.wId2, &r.wAbs, &r.wLabel, &r.isDef, &r.hasBind, &r.cAt, &r.uAt)
	if err != nil {
		return nil, err
	}

	return r.toWorkDirectory(), nil
}

func queryDirectWorkDirectory(db *DB, targetPath string) (*model.WorkDirectory, error) {
	row := db.conn.QueryRow(constants.SQLResolveWorkDirDirectPrefix, targetPath, targetPath, targetPath)
	item, err := scanWorkDirectoryRow(row)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func queryDefaultWorkDirectory(db *DB) (*model.WorkDirectory, error) {
	row := db.conn.QueryRow(constants.SQLSelectDefaultWorkDirectory)
	item, err := scanWorkDirectoryRow(row)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func queryFallbackWorkDirectory(db *DB, targetPath string) (*model.WorkDirectory, error) {
	direct, errDirect := queryDirectWorkDirectory(db, targetPath)
	if errDirect == nil && direct != nil {
		return direct, nil
	}

	return queryDefaultWorkDirectory(db)
}

// ResolveWorkDirectoryForPath finds the effective WorkDirectory using hierarchical longest-prefix matching.
func (db *DB) ResolveWorkDirectoryForPath(targetPath string) (*model.WorkDirectory, error) {
	matched, errNearest := queryNearestFolderWorkDir(db, targetPath)
	if errNearest == nil && matched != nil {
		return matched, nil
	}

	fallback, errFallback := queryFallbackWorkDirectory(db, targetPath)
	if errFallback != nil {
		return nil, nil
	}

	return fallback, nil
}

func (db *DB) getFolderGroupItem(groupName, folderPath string) (*model.FolderGroup, error) {
	const sqlFind = "SELECT FolderGroupId, GroupName, FolderPath, Description, IsDefault, IsActive, CreatedAt, UpdatedAt FROM FolderGroup WHERE GroupName = ? COLLATE NOCASE AND FolderPath = ? COLLATE NOCASE"
	row := db.conn.QueryRow(sqlFind, groupName, folderPath)
	item, err := scanFolderGroupRow(row)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// AddFolderToGroup adds a folder path to a named folder group.
func (db *DB) AddFolderToGroup(groupName, folderPath, description string) (*model.FolderGroup, error) {
	_, err := db.conn.Exec(constants.SQLUpsertFolderGroup, groupName, folderPath, description, 0, 1)
	if err != nil {
		return nil, apperror.WrapSimple(err, "AddFolderToGroup.Exec")
	}

	return db.getFolderGroupItem(groupName, folderPath)
}

// RemoveFolderFromGroup removes a folder path from a named group.
func (db *DB) RemoveFolderFromGroup(groupName, folderPath string) error {
	_, err := db.conn.Exec(constants.SQLDeleteFolderFromGroup, groupName, folderPath)
	if err != nil {
		return apperror.WrapSimple(err, "RemoveFolderFromGroup")
	}

	return nil
}

// ListFoldersInGroup lists all folders in the specified group.
func (db *DB) ListFoldersInGroup(groupName string) ([]model.FolderGroup, error) {
	rows, err := db.conn.Query(constants.SQLSelectFoldersInGroup, groupName)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListFoldersInGroup.Query")
	}

	defer rows.Close()

	return collectFolderGroups(rows)
}

func collectFolderGroups(rows *sql.Rows) ([]model.FolderGroup, error) {
	var list []model.FolderGroup
	for rows.Next() {
		item, err := scanFolderGroupRow(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *item)
	}

	return list, nil
}

// ListActiveGroups lists all active folder groups with aggregate statistics.
func (db *DB) ListActiveGroups() ([]model.FolderGroupSummary, error) {
	rows, err := db.conn.Query(constants.SQLSelectActiveGroupsSummary)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListActiveGroups.Query")
	}

	defer rows.Close()

	return collectGroupSummaries(rows)
}

func collectGroupSummaries(rows *sql.Rows) ([]model.FolderGroupSummary, error) {
	var list []model.FolderGroupSummary
	for rows.Next() {
		var s model.FolderGroupSummary
		if err := rows.Scan(&s.GroupName, &s.FolderCount, &s.CreatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "collectGroupSummaries.Scan")
		}

		list = append(list, s)
	}

	return list, nil
}

func (db *DB) getSshKeyAssocByTarget(assocType, targetId string, sshKeyId int64) (*model.SshKeyAssociation, error) {
	const sqlFind = `SELECT SshKeyAssociationId, SshKeyId, AssociationType, TargetIdentifier, FolderPath, GroupName,
		WorkDirectoryId, RepoId, GitAuthorName, GitAuthorEmail, IsActive, IsAutoInjected, Priority, CreatedAt, UpdatedAt
		FROM SshKeyAssociation WHERE AssociationType = ? AND TargetIdentifier = ? COLLATE NOCASE AND SshKeyId = ?`
	row := db.conn.QueryRow(sqlFind, assocType, targetId, sshKeyId)
	item, err := scanSshKeyAssociationRow(row)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// AssociateSSHKeyWithFolder associates an SSH key with a folder path and optional author overrides.
func (db *DB) AssociateSSHKeyWithFolder(sshKeyId int64, folderPath, authorName, authorEmail string) (*model.SshKeyAssociation, error) {
	_, err := db.conn.Exec(constants.SQLUpsertSshKeyAssociation, sshKeyId, constants.AssocTypeFolder, folderPath, folderPath, "", nil, nil, authorName, authorEmail, 1, 1, 100)
	if err != nil {
		return nil, apperror.WrapSimple(err, "AssociateSSHKeyWithFolder")
	}

	return db.getSshKeyAssocByTarget(constants.AssocTypeFolder, folderPath, sshKeyId)
}

// AssociateSSHKeyWithGroup associates an SSH key with a named workgroup and optional author overrides.
func (db *DB) AssociateSSHKeyWithGroup(sshKeyId int64, groupName, authorName, authorEmail string) (*model.SshKeyAssociation, error) {
	_, err := db.conn.Exec(constants.SQLUpsertSshKeyAssociation, sshKeyId, constants.AssocTypeGroup, groupName, "", groupName, nil, nil, authorName, authorEmail, 1, 1, 100)
	if err != nil {
		return nil, apperror.WrapSimple(err, "AssociateSSHKeyWithGroup")
	}

	return db.getSshKeyAssocByTarget(constants.AssocTypeGroup, groupName, sshKeyId)
}

// AssociateSSHKeyWithWorkDir associates an SSH key with a specific WorkDirectory.
func (db *DB) AssociateSSHKeyWithWorkDir(sshKeyId int64, workDirId int64) (*model.SshKeyAssociation, error) {
	targetId := strconv.FormatInt(workDirId, 10)
	_, err := db.conn.Exec(constants.SQLUpsertSshKeyAssociation, sshKeyId, constants.AssocTypeWorkDir, targetId, "", "", workDirId, nil, "", "", 1, 1, 100)
	if err != nil {
		return nil, apperror.WrapSimple(err, "AssociateSSHKeyWithWorkDir")
	}

	return db.getSshKeyAssocByTarget(constants.AssocTypeWorkDir, targetId, sshKeyId)
}

// AssociateSSHKeyWithRepo associates an SSH key with a specific repository.
func (db *DB) AssociateSSHKeyWithRepo(sshKeyId int64, repoId int64) (*model.SshKeyAssociation, error) {
	targetId := strconv.FormatInt(repoId, 10)
	_, err := db.conn.Exec(constants.SQLUpsertSshKeyAssociation, sshKeyId, constants.AssocTypeRepo, targetId, "", "", nil, repoId, "", "", 1, 1, 100)
	if err != nil {
		return nil, apperror.WrapSimple(err, "AssociateSSHKeyWithRepo")
	}

	return db.getSshKeyAssocByTarget(constants.AssocTypeRepo, targetId, sshKeyId)
}

// RemoveSSHKeyAssociation removes an SSH key association record.
func (db *DB) RemoveSSHKeyAssociation(assocId int64) error {
	_, err := db.conn.Exec(constants.SQLDeleteSshKeyAssociation, assocId)
	if err != nil {
		return apperror.WrapSimple(err, "RemoveSSHKeyAssociation")
	}

	return nil
}

// ListSSHKeyAssociations returns all SSH key associations.
func (db *DB) ListSSHKeyAssociations() ([]model.SshKeyAssociation, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllSshKeyAssociations)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListSSHKeyAssociations.Query")
	}

	defer rows.Close()

	return collectSshKeyAssociations(rows)
}

func collectSshKeyAssociations(rows *sql.Rows) ([]model.SshKeyAssociation, error) {
	var list []model.SshKeyAssociation
	for rows.Next() {
		item, err := scanSshKeyAssociationRow(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *item)
	}

	return list, nil
}

func resolveWorkDirIdForPath(db *DB, repoPath string) int64 {
	wd, err := db.ResolveWorkDirectoryForPath(repoPath)
	if err != nil || wd == nil {
		return 0
	}

	return wd.WorkDirectoryId
}

type effectiveCredScanned struct {
	assocId, keyId, prio                                    int64
	aType, tId, kName, pPath, pubKey, fPrint, aName, aEmail string
	isAuto                                                  int
}

func (s *effectiveCredScanned) toModel() *model.EffectiveFolderCredentials {
	return &model.EffectiveFolderCredentials{
		HasCustomKey:    true,
		SshKeyId:        s.keyId,
		KeyName:         s.kName,
		PrivatePath:     s.pPath,
		PublicKey:       s.pubKey,
		Fingerprint:     s.fPrint,
		AssociationType: s.aType,
		GitAuthorName:   s.aName,
		GitAuthorEmail:  s.aEmail,
		IsAutoInjected:  (s.isAuto == 1),
	}
}

func scanEffectiveCredRow(row *sql.Row) (*model.EffectiveFolderCredentials, error) {
	var s effectiveCredScanned
	err := row.Scan(&s.assocId, &s.aType, &s.tId, &s.keyId, &s.kName, &s.pPath, &s.pubKey, &s.fPrint, &s.aName, &s.aEmail, &s.isAuto, &s.prio)
	if err != nil {
		return nil, err
	}

	return s.toModel(), nil
}

func buildDefaultCred(keyId int64, kName, pPath, pubKey, fPrint string) *model.EffectiveFolderCredentials {
	return &model.EffectiveFolderCredentials{
		HasCustomKey:    true,
		SshKeyId:        keyId,
		KeyName:         kName,
		PrivatePath:     pPath,
		PublicKey:       pubKey,
		Fingerprint:     fPrint,
		AssociationType: constants.AssocTypeDefault,
		IsAutoInjected:  true,
	}
}

func queryDefaultSshKeyCred(db *DB) (*model.EffectiveFolderCredentials, error) {
	row := db.conn.QueryRow(constants.SQLSelectDefaultSSHKeySimple)
	var keyId int64
	var kName, pPath, pubKey, fPrint string

	if err := row.Scan(&keyId, &kName, &pPath, &pubKey, &fPrint); err != nil {
		return &model.EffectiveFolderCredentials{HasCustomKey: false}, nil
	}

	return buildDefaultCred(keyId, kName, pPath, pubKey, fPrint), nil
}

// ResolveEffectiveCredentialsForPath resolves the active SSH key and author overrides via 5-stage cascading priority.
func (db *DB) ResolveEffectiveCredentialsForPath(repoPath string, repoId int64) (*model.EffectiveFolderCredentials, error) {
	workDirId := resolveWorkDirIdForPath(db, repoPath)
	row := db.conn.QueryRow(constants.SQLResolveEffectiveCredentials, repoId, repoPath, repoPath, repoPath, repoPath, repoPath, repoPath, workDirId)

	cred, err := scanEffectiveCredRow(row)
	if err == nil && cred != nil {
		return cred, nil
	}

	return queryDefaultSshKeyCred(db)
}
