import re

with open('cli/cmd/repo_create_init.go', 'r') as f:
    content = f.read()

content = content.replace('''
	isExisting := false
	if _, statErr := os.Stat(absDir); statErr == nil {
		isExisting = true
	} else {
		if mkErr := os.MkdirAll(absDir, 0755); mkErr != nil {
			return apperror.WrapSimple(mkErr, "create directory:")
		}
	}
''', '''
	isExisting := checkDirExists(absDir)
	if !isExisting {
		err := os.MkdirAll(absDir, 0755)
		if err != nil {
			return apperror.WrapSimple(err, "create directory:")
		}
	}
''')

content = content.replace('''
	gitDir := filepath.Join(absDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmdInit := exec.Command("git", "init", "-b", "main")
		cmdInit.Dir = absDir
		if initErr := cmdInit.Run(); initErr != nil {
			return apperror.WrapSimple(initErr, "git init:")
		}
	}
''', '''
	gitDir := filepath.Join(absDir, ".git")
	_, err := os.Stat(gitDir)
	if os.IsNotExist(err) {
		errInit := runGitInit(absDir)
		if errInit != nil {
			return errInit
		}
	}
''')

content = content.replace('''
	if !isExisting {
		writeInitialFiles(absDir, p)
		if commitErr := commitInitialFiles(absDir); commitErr != nil {
			return commitErr
		}
	}

	if p.IsCommon {
		if err := applyCommonFiles(absDir); err != nil {
			return err
		}
	}

	if p.IsCG {
		if err := applyCGFiles(absDir); err != nil {
			return err
		}
	}

	return nil
}
''', '''
	if !isExisting {
		writeInitialFiles(absDir, p)
		commitErr := commitInitialFiles(absDir)
		if commitErr != nil {
			return commitErr
		}
	}

	errCommon := maybeApplyCommonFiles(p.IsCommon, absDir)
	if errCommon != nil {
		return errCommon
	}

	errCG := maybeApplyCGFiles(p.IsCG, absDir)
	if errCG != nil {
		return errCG
	}

	return nil
}

func checkDirExists(absDir string) bool {
	_, statErr := os.Stat(absDir)
	return statErr == nil
}

func runGitInit(absDir string) error {
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = absDir
	initErr := cmdInit.Run()
	if initErr != nil {
		return apperror.WrapSimple(initErr, "git init:")
	}
	return nil
}

func maybeApplyCommonFiles(isCommon bool, absDir string) error {
	if !isCommon {
		return nil
	}
	return applyCommonFiles(absDir)
}

func maybeApplyCGFiles(isCG bool, absDir string) error {
	if !isCG {
		return nil
	}
	return applyCGFiles(absDir)
}
''')

with open('cli/cmd/repo_create_init.go', 'w') as f:
    f.write(content)
