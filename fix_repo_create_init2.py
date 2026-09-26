import re

with open('cli/cmd/repo_create_init.go', 'r') as f:
    content = f.read()

content = content.replace('''
	isExisting := checkDirExists(absDir)
	if !isExisting {
		err := os.MkdirAll(absDir, 0755)
		if err != nil {
			return apperror.WrapSimple(err, "create directory:")
		}
	}
''', '''
	isExisting := checkDirExists(absDir)
	errMk := maybeMkdirAll(isExisting, absDir)
	if errMk != nil {
		return errMk
	}
''')

content = content.replace('''
	gitDir := filepath.Join(absDir, ".git")
	_, err := os.Stat(gitDir)
	if os.IsNotExist(err) {
		errInit := runGitInit(absDir)
		if errInit != nil {
			return errInit
		}
	}
''', '''
	gitDir := filepath.Join(absDir, ".git")
	errInit := maybeGitInit(gitDir, absDir)
	if errInit != nil {
		return errInit
	}
''')

content = content.replace('''
	if !isExisting {
		writeInitialFiles(absDir, p)
		commitErr := commitInitialFiles(absDir)
		if commitErr != nil {
			return commitErr
		}
	}
''', '''
	errInitFiles := maybeWriteInitialFiles(isExisting, absDir, p)
	if errInitFiles != nil {
		return errInitFiles
	}
''')

content = content.replace('''
func checkDirExists(absDir string) bool {
''', '''
func maybeMkdirAll(isExisting bool, absDir string) error {
	if isExisting {
		return nil
	}
	err := os.MkdirAll(absDir, 0755)
	if err != nil {
		return apperror.WrapSimple(err, "create directory:")
	}
	return nil
}

func maybeGitInit(gitDir, absDir string) error {
	_, err := os.Stat(gitDir)
	if !os.IsNotExist(err) {
		return nil
	}
	return runGitInit(absDir)
}

func maybeWriteInitialFiles(isExisting bool, absDir string, p createRepoParams) error {
	if isExisting {
		return nil
	}
	writeInitialFiles(absDir, p)
	return commitInitialFiles(absDir)
}

func checkDirExists(absDir string) bool {
''')

with open('cli/cmd/repo_create_init.go', 'w') as f:
    f.write(content)
