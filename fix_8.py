import re

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

r('gitmap/cmd/githubdesktop.go', 'output, runErr :=', '_, runErr :=')
r('gitmap/cmd/importcmd.go', 'return apperror.Wrap(err, constants.MsgImportFailed, nil)', 'panic(err)')
r('gitmap/cmd/install.go', 'return apperror.New("install-deps requires an argument. E.g. gitmap install-deps all", "E9000", nil)', 'panic("install-deps requires an argument. E.g. gitmap install-deps all")')
r('gitmap/cmd/install_custom.go', 'apperror.Wrap(tool, "Error unknown tool", nil)', 'apperror.New("Error unknown tool: " + tool, "E9000", nil)')
r('gitmap/cmd/installcleancode.go', 'apperror.Wrap(constants.DefaultCleanCodeURL, "Failed to download", nil)', 'apperror.New("Failed to download", "E9000", nil)')
r('gitmap/cmd/installtools.go', 'return apperror.Wrap(err, "Download failed", nil)', 'panic(err)')
r('gitmap/cmd/latestbranch.go', 'return apperror.Wrap(cfg.remote, "constants.ErrLatestBranchNoRefs", nil)', 'return nil')
r('gitmap/cmd/latestbranch.go', 'return apperror.New("fatal error", "E9000", nil)', 'panic("fatal")')
