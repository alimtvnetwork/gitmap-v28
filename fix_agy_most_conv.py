import re

with open('cli/cmdagy/agy_most_conv.go', 'r') as f:
    content = f.read()

content = content.replace('''
		if len(args) > 0 {
			if val, err := strconv.Atoi(args[0]); err == nil && val > 0 {
				n = val
			}
		}
''', '''
		n = parseMostConvArgCount(args, n)
''')

content = content.replace('''
func runMostConv(n int) error {
''', '''
func parseMostConvArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}

func runMostConv(n int) error {
''')

with open('cli/cmdagy/agy_most_conv.go', 'w') as f:
    f.write(content)
