import sys

def process(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # fix applyEnvDelete
    content = content.replace("fmt.Printf(constants.MsgEnvDeleted, name)\n}", "fmt.Printf(constants.MsgEnvDeleted, name)\n\treturn nil\n}")

    # fix applyEnvPathAdd
    content = content.replace("fmt.Printf(constants.MsgEnvDryPath, dir)\n\t\treturn\n\t}", "fmt.Printf(constants.MsgEnvDryPath, dir)\n\t\treturn nil\n\t}")
    content = content.replace("fmt.Printf(constants.MsgEnvPathAdded, dir)\n}", "fmt.Printf(constants.MsgEnvPathAdded, dir)\n\treturn nil\n}")

    # fix applyEnvPathRemove
    content = content.replace("fmt.Printf(constants.MsgEnvDryDelete, dir)\n\t\treturn\n\t}", "fmt.Printf(constants.MsgEnvDryDelete, dir)\n\t\treturn nil\n\t}")
    content = content.replace("fmt.Printf(constants.MsgEnvPathRemoved, dir)\n}", "fmt.Printf(constants.MsgEnvPathRemoved, dir)\n\treturn nil\n}")

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

process('gitmap/cmd/envops.go')
