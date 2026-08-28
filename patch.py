import sys

def process(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    if "envplatform_windows.go" in filepath:
        content = content.replace("func setEnvPersistent(name, value string, system bool, _ string) {", "func setEnvPersistent(name, value string, system bool, _ string) error {")
        content = content.replace("runSetx(args)", "return runSetx(args)")
        
        content = content.replace("func deleteEnvPersistent(name string, system bool, _ string) {", "func deleteEnvPersistent(name string, system bool, _ string) error {")
        
        content = content.replace("func addPathPersistent(dir string, system bool, _ string) {", "func addPathPersistent(dir string, system bool, _ string) error {")
        content = content.replace('setEnvPersistent("PATH", newPath, system, "")', 'return setEnvPersistent("PATH", newPath, system, "")')
        
        content = content.replace("func removePathPersistent(dir string, system bool, _ string) {", "func removePathPersistent(dir string, system bool, _ string) error {")
        
    elif "envplatform_unix.go" in filepath:
        content = content.replace("func setEnvPersistent(name, value string, _ bool, shell string) {", "func setEnvPersistent(name, value string, _ bool, shell string) error {")
        content = content.replace("appendToProfile(profilePath, name, exportLine)", "return appendToProfile(profilePath, name, exportLine)")
        
        content = content.replace("func deleteEnvPersistent(name string, _ bool, shell string) {", "func deleteEnvPersistent(name string, _ bool, shell string) error {")
        content = content.replace("removeFromProfile(profilePath, name)", "return removeFromProfile(profilePath, name)")
        
        content = content.replace("func addPathPersistent(dir string, _ bool, shell string) {", "func addPathPersistent(dir string, _ bool, shell string) error {")
        content = content.replace("appendLineToProfile(profilePath, exportLine, marker)", "return appendLineToProfile(profilePath, exportLine, marker)")
        
        content = content.replace("func removePathPersistent(dir string, _ bool, shell string) {", "func removePathPersistent(dir string, _ bool, shell string) error {")
        content = content.replace("removeLineFromProfile(profilePath, marker)", "return removeLineFromProfile(profilePath, marker)")
        
        content = content.replace("func appendToProfile(profilePath, name, exportLine string) {", "func appendToProfile(profilePath, name, exportLine string) error {")
        content = content.replace("writeProfileContent(profilePath, updatedContent)", "return writeProfileContent(profilePath, updatedContent)")
        
        content = content.replace("func removeFromProfile(profilePath, name string) {", "func removeFromProfile(profilePath, name string) error {")
        
        content = content.replace("func appendLineToProfile(profilePath, line, marker string) {", "func appendLineToProfile(profilePath, line, marker string) error {")
        
        content = content.replace("func removeLineFromProfile(profilePath, marker string) {", "func removeLineFromProfile(profilePath, marker string) error {")
        content = content.replace("writeProfileContent(profilePath, strings.Join(filtered, \"\\n\"))", "return writeProfileContent(profilePath, strings.Join(filtered, \"\\n\"))")
        
        content = content.replace("func writeProfileContent(path, content string) {", "func writeProfileContent(path, content string) error {")
        content = content.replace("return apperror.NewSimple(constants.ErrEnvProfileWrite, \"E9000\")\n\t}", "return apperror.NewSimple(constants.ErrEnvProfileWrite, \"E9000\")\n\t}\n\n\treturn nil")

    elif "envops.go" in filepath:
        content = content.replace("func applyEnvSet(name, value string, f envSetFlags) {", "func applyEnvSet(name, value string, f envSetFlags) error {")
        content = content.replace("fmt.Printf(constants.MsgEnvDrySet, name, value)\n\t\treturn", "fmt.Printf(constants.MsgEnvDrySet, name, value)\n\t\treturn nil")
        content = content.replace("setEnvPersistent(name, value, f.system, f.shell)", "if err := setEnvPersistent(name, value, f.system, f.shell); err != nil {\n\t\treturn err\n\t}")
        content = content.replace("fmt.Printf(constants.MsgEnvSet, name, value)\n}", "fmt.Printf(constants.MsgEnvSet, name, value)\n\treturn nil\n}")
        
        content = content.replace("applyEnvSet(name, value, flags)\n\treturn nil", "return applyEnvSet(name, value, flags)")
        
        content = content.replace("func applyEnvDelete(name string, f envCommonFlags) {", "func applyEnvDelete(name string, f envCommonFlags) error {")
        content = content.replace("fmt.Printf(constants.MsgEnvDryDelete, name)\n\t\treturn", "fmt.Printf(constants.MsgEnvDryDelete, name)\n\t\treturn nil")
        content = content.replace("deleteEnvPersistent(name, f.system, f.shell)", "if err := deleteEnvPersistent(name, f.system, f.shell); err != nil {\n\t\treturn err\n\t}")
        content = content.replace("fmt.Printf(constants.MsgEnvDelete, name)\n}", "fmt.Printf(constants.MsgEnvDelete, name)\n\treturn nil\n}")
        
        content = content.replace("applyEnvDelete(name, flags)\n\treturn nil", "return applyEnvDelete(name, flags)")
        
        content = content.replace("func applyEnvPathAdd(dir string, f envCommonFlags) {", "func applyEnvPathAdd(dir string, f envCommonFlags) error {")
        content = content.replace("fmt.Printf(constants.MsgEnvDryPathAdd, dir)\n\t\treturn", "fmt.Printf(constants.MsgEnvDryPathAdd, dir)\n\t\treturn nil")
        content = content.replace("addPathPersistent(dir, f.system, f.shell)", "if err := addPathPersistent(dir, f.system, f.shell); err != nil {\n\t\treturn err\n\t}")
        content = content.replace("fmt.Printf(constants.MsgEnvPathAdd, dir)\n}", "fmt.Printf(constants.MsgEnvPathAdd, dir)\n\treturn nil\n}")
        
        content = content.replace("applyEnvPathAdd(dir, flags)\n\treturn nil", "return applyEnvPathAdd(dir, flags)")
        
        content = content.replace("func applyEnvPathRemove(dir string, f envCommonFlags) {", "func applyEnvPathRemove(dir string, f envCommonFlags) error {")
        content = content.replace("fmt.Printf(constants.MsgEnvDryPathRemove, dir)\n\t\treturn", "fmt.Printf(constants.MsgEnvDryPathRemove, dir)\n\t\treturn nil")
        content = content.replace("removePathPersistent(dir, f.system, f.shell)", "if err := removePathPersistent(dir, f.system, f.shell); err != nil {\n\t\treturn err\n\t}")
        content = content.replace("fmt.Printf(constants.MsgEnvPathRemove, dir)\n}", "fmt.Printf(constants.MsgEnvPathRemove, dir)\n\treturn nil\n}")
        
        content = content.replace("applyEnvPathRemove(dir, flags)\n\treturn nil", "return applyEnvPathRemove(dir, flags)")

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

process('gitmap/cmd/envplatform_windows.go')
process('gitmap/cmd/envplatform_unix.go')
process('gitmap/cmd/envops.go')
