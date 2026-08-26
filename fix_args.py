import re

path = "gitmap/cmd/agy_cmd.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

repl = """func dispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && (args[0] == "agy" || args[0] == "ag" || args[0] == "antigravity") {
		args = args[1:]
	}
	AgyCmd.SetArgs(args)
	return AgyCmd.ExecuteContext(ctx)
}"""

content = re.sub(r'func dispatchAgy\(ctx context.Context, args \[\]string, root \*cobra.Command\) error \{[\s\S]*?return AgyCmd.ExecuteContext\(ctx\)\n\}', repl, content)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
