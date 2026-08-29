# Task: Fix Inverted Booleans

According to the Coding Guidelines: 'Always use explicit boolean state variables (e.g., `isFail`) instead of inverting positive ones (`!isSuccess`).'

## Files to fix

- d:\wp-work\riseup-asia\gitmap\gitmap\clonefrom\jsonschema_test.go:172: isNotNumber := !isNumber
- d:\wp-work\riseup-asia\gitmap\gitmap\cluster\distribution.go:28: hasNoClients := len(clients) == EmptySize
- d:\wp-work\riseup-asia\gitmap\gitmap\cluster\exec_proj.go:62: isNotFound := foundPath == emptyString
- d:\wp-work\riseup-asia\gitmap\gitmap\cluster\node_resolver.go:41: hasNoNodes := len(allNodes) == constants.EmptySliceLength
- d:\wp-work\riseup-asia\gitmap\gitmap\cluster\ui.go:56: hasNoClients := numClients == DefaultProgress
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\auditlegacy.go:82: isNotAuditScannable := !isAuditScannable(path)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\cdops.go:174: hasNoGroup := len(group) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\chromeprofile.go:96: isNotRunning := !isRunning
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\chromeprofile_merge.go:57: isNotKnownMergeWhat := !isKnownMergeWhat(*what)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\chrome_bookmarks.go:75: hasNoRootsAfterRootName := rootName != "" && len(roots) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\chrome_bookmarks.go:83: hasNoRootsAfterFolderPath := folderPath != "" && len(roots) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\chrome_bookmarks.go:92: hasNoRootsAfterMatch := hasMatchOrTitle == true && len(roots) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clone.go:573: isNotDirectURL := !isDirectURL(in)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonefixrepo_modifiers.go:48: isNotCfrModifierToken := !isCfrModifierToken(args[i])
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonenextfolderdispatch.go:65: isNotFolderShaped := !isFolderShaped(token)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonenextfolderdispatch.go:106: isNotFolderShaped := !isFolderShaped(second)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonenextfolderdispatch_test.go:128: isNotFolderShaped := !isFolderShaped(h)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonepmsync.go:90: isNotOn := !isOn
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clonevscode.go:28: isNotVSCodeAvailable := !isVSCodeAvailable()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clone_ssh_hostkey.go:11: isNotAssumeYes := !isAssumeYes
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clone_stale_binary_test.go:52: isNotDirectURL := !isDirectURL(u)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\cluster.go:18: hasNoArgs := len(args) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clustercommand.go:32: hasNoSubCmds := len(subCmds) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\clustercommand.go:66: hasNoEffectiveNodes := len(effective) == 0
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\fixauth.go:37: isNotGitRepoCWD := !isGitRepoCWD()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\githubdesktop.go:27: isNotGitRepo := !isGitRepo(target)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\gomodreplace_test.go:38: isNotExcludedDir := !isExcludedDir(".git")
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\gomodreplace_test.go:45: isNotExcludedDir := !isExcludedDir("vendor")
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\gomodreplace_test.go:52: isNotExcludedDir := !isExcludedDir("node_modules")
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\gomod_integration_test.go:132: isNotWorkTreeDirty := !isWorkTreeDirty()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\hasanyupdates.go:17: isNotInsideGitRepo := !isInsideGitRepo()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\installctx_linux_e2e_test.go:244: isNotSingleBlock := count != 1
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\pull.go:155: isNotGitRepoCWD := !isGitRepoCWD()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\pullreleasecd.go:135: isNotPRCURL := !isPRCURL(token)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\push.go:135: isNotGitRepoCWD := !isGitRepoCWD()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\push.go:220: isNotNonFastForwardRejection := !isNonFastForwardRejection(stderr)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\reclone_confirm.go:49: isNotStdinInteractive := !isStdinInteractive()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\reclone_validate.go:174: isNotSchemeChar := !isSchemeChar(ch)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\regoldens_diff.go:42: isNotGitWorkingTree := !isGitWorkingTree()
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\regoldens_diff.go:101: isNotGoldenFixturePath := !isGoldenFixturePath(path)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\regoldens_diff.go:159: isNotGoldenFixturePath := !isGoldenFixturePath(fields[2])
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\regoldens_diff_test.go:110: isNotGoldenFixturePath := !isGoldenFixturePath(p)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\releasealias_git.go:32: isNotWorkingTreeDirty := !isWorkingTreeDirty(target)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\replacewalk_test.go:15: isNotExcludedDir := !isExcludedDir(name)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\replacewalk_test.go:31: isNotExcludedRelease := !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release"))
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\replacewalk_test.go:35: isNotExcludedV1 := !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release", "v1"))
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\replacewalk_test.go:39: isNotExcludedAssets := !isExcludedPrefix(root, filepath.Join(root, ".gitmap", "release-assets"))
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\replacewalk_test.go:93: isNotBinaryFile := !isBinaryFile(binPath)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\reporeclone.go:74: isNotGitRepoDir := !isGitRepoDir(cwd)
- d:\wp-work\riseup-asia\gitmap\gitmap\cmd\reporeclone.go:91: isNotGitRepoDir := !isGitRepoDir(abs)
