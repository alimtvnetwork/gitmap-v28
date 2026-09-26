import re

with open('cli/cmdagy/agy_conv_ls.go', 'r') as f:
    content = f.read()

content = content.replace('''
		if len(args) > 0 {
			if val, err := strconv.Atoi(args[0]); err == nil && val > 0 {
				n = val
			}
		}
''', '''
		n = parseAgyConvLsArgCount(args, n)
''')

content = content.replace('''
func init() {
''', '''
func parseAgyConvLsArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}

func init() {
''')

with open('cli/cmdagy/agy_conv_ls.go', 'w') as f:
    f.write(content)
