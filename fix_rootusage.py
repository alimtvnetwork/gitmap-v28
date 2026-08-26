import re

path = "gitmap/cmd/rootusage.go"
with open(path, "r", encoding="utf-8") as f:
    content = f.read()

# I want to add an "AUTOMATION & INSTALLERS" and "INTEGRATIONS" super category.
# I'll put them in printUsageAdvancedCategories()
repl = """// printUsageAdvancedCategories prints the Cluster and Advanced categories.
func printUsageAdvancedCategories() {
	printSuperCategory("INSTALL & AUTOMATION", func() {
		printGroupInstallers()
	})
	printSuperCategory("INTEGRATIONS", func() {
		printGroupIntegrations()
	})
	printSuperCategory("CLUSTER & NETWORK", func() {
"""
content = content.replace('// printUsageAdvancedCategories prints the Cluster and Advanced categories.\nfunc printUsageAdvancedCategories() {\n\tprintSuperCategory("CLUSTER & NETWORK", func() {', repl)

with open(path, "w", encoding="utf-8") as f:
    f.write(content)

print("done")
