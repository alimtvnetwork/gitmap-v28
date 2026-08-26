import re, glob

for file in glob.glob("gitmap/cmd/rootusage*.go"):
    with open(file, "r", encoding="utf-8") as f:
        content = f.read()
    
    # We want to replace `fmt.Println(colorGroupHeader(SOMETHING))` 
    # with `fmt.Println(colorGroupHeader(SOMETHING))\n\tfmt.Println()`
    # but only if it's not already followed by a newline print.
    
    def repl(m):
        header_print = m.group(0)
        # Check if the next line is already fmt.Println()
        return header_print + "\n\tfmt.Println()"
    
    # Simple regex: find fmt.Println(colorGroupHeader(XXX))
    # Using a negative lookahead to avoid adding it twice
    content = re.sub(r'(\tfmt\.Println\(colorGroupHeader\([^)]+\)\))(?!\s*fmt\.Println\(\))', repl, content)
    
    with open(file, "w", encoding="utf-8") as f:
        f.write(content)
print("done")
