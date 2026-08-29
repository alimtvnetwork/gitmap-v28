import os
import re

gitmap_dir = r"d:\work\gitmap\gitmap"

# 1. Revert go.mod
mod_path = os.path.join(gitmap_dir, "go.mod")
with open(mod_path, "r", encoding="utf-8") as f: c = f.read()
c = c.replace("go 1.25.0", "go 1.24.0")
with open(mod_path, "w", encoding="utf-8") as f: f.write(c)

# 2. Fix Spellings and "GET"
replacements = {
    "recognised": "recognized",
    "recognise": "recognize",
    "initialises": "initializes",
    "honours": "honors",
    "honouring": "honoring",
    "honoured": "honored",
    "behaviour": "behavior",
    "recognisable": "recognizable",
    "unrecognised": "unrecognized",
    "cancelled": "canceled",
    "Centralised": "Centralized",
    "labelled": "labeled",
    "marshalled": "marshaled",
    '"GET"': "http.MethodGet"
}

def fix_file(filepath):
    with open(filepath, "r", encoding="utf-8") as f: content = f.read()
    orig = content
    for old, new in replacements.items():
        content = content.replace(old, new)
        
    # Also fix gosimple == true/false safely where possible
    content = re.sub(r'(\w+)\s*==\s*true', r'\1', content)
    content = re.sub(r'(\w+)\s*==\s*false', r'!\1', content)
        
    if content != orig:
        with open(filepath, "w", encoding="utf-8", newline="\n") as f:
            f.write(content)

for root, _, files in os.walk(gitmap_dir):
    for file in files:
        if file.endswith(".go"):
            fix_file(os.path.join(root, file))

# 3. Unreachable Code
def del_unreachable_return(filepath, report_line):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: lines = f.readlines()
    
    if "return nil" in lines[report_line-1]:
        lines.pop(report_line-1)
    else:
        found = False
        for i in range(report_line-2, max(-1, report_line-6), -1):
            if "return nil" in lines[i]:
                lines.pop(i)
                found = True
                break
            
    with open(path, "w", encoding="utf-8") as f: f.writelines(lines)

del_unreachable_return(r"cmd\asops.go", 54)
del_unreachable_return(r"cmd\branch.go", 36)
del_unreachable_return(r"cmd\changelog.go", 58) # higher line first
del_unreachable_return(r"cmd\changelog.go", 36)
del_unreachable_return(r"cmd\cluster.go", 56)
del_unreachable_return(r"cmd\dbreset.go", 22)
del_unreachable_return(r"cmd\historyreset.go", 24)
del_unreachable_return(r"cmd\probe.go", 61)
del_unreachable_return(r"cmd\regoldens.go", 120)
del_unreachable_return(r"cmd\rm.go", 83)
del_unreachable_return(r"cmd\selfuninstallhandoff.go", 97)
del_unreachable_return(r"cmd\tasksync.go", 72)
del_unreachable_return(r"cmd\watch.go", 88)

def del_panic_return(filepath, search_str, after_str):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: lines = f.readlines()
    for i, l in enumerate(lines):
        if search_str in l and i > 0 and after_str in lines[i-1]:
            lines.pop(i)
            break
        if search_str in l and i > 1 and after_str in lines[i-2]:
            lines.pop(i)
            break
    with open(path, "w", encoding="utf-8") as f: f.writelines(lines)

del_panic_return(r"cmd\reinstall.go", "return constants.ReinstallModeSelf", "panic(")
del_panic_return(r"cmd\replaceflags.go", "return true", "panic(")
del_panic_return(r"cmd\taskops.go", "return model.TaskEntry{}", "panic(")
del_panic_return(r"cmd\update.go", 'return ""', "panic(")

# 4. Errorlint Fixes
def replace_exact(filepath, old, new):
    path = os.path.join(gitmap_dir, filepath)
    with open(path, "r", encoding="utf-8") as f: c = f.read()
    c = c.replace(old, new)
    with open(path, "w", encoding="utf-8") as f: f.write(c)

replace_exact(r"cluster\exec_cmd.go", "if exitErr, ok := err.(*exec.ExitError); ok {", "var exitErr *exec.ExitError\n\tif errors.As(err, &exitErr) {")
replace_exact(r"cluster\exec_lifecycle.go", "if exitErr, ok := err.(*exec.ExitError); ok {", "var exitErr *exec.ExitError\n\tif errors.As(err, &exitErr) {")
replace_exact(r"cluster\exec_ps.go", "exitErr, isExitErr := err.(*exec.ExitError)", "var exitErr *exec.ExitError\n\tisExitErr := errors.As(err, &exitErr)")
replace_exact(r"cmd\chrome_backup.go", "if err == io.EOF {", "if errors.Is(err, io.EOF) {")
replace_exact(r"cmd\helpdashboard_download.go", "err != io.EOF", "!errors.Is(err, io.EOF)")
replace_exact(r"cmd\root.go", "if appErr, ok := err.(*apperror.AppError); ok {", "var appErr *apperror.AppError\n\t\tif errors.As(err, &appErr) {")
replace_exact(r"cmd\root.go", "if appErr, ok := err.(*apperror.AppError); ok && appErr.Cause != nil {", "var appErr *apperror.AppError\n\t\tif errors.As(err, &appErr) && appErr.Cause != nil {")
replace_exact(r"cmd\chrome_backup_integration_test.go", "if err == io.EOF {", "if errors.Is(err, io.EOF) {")
replace_exact(r"store\ssh_repo_test.go", "if appErr.Cause != apperror.ErrNotFound {", "if !errors.Is(appErr.Cause, apperror.ErrNotFound) {")

