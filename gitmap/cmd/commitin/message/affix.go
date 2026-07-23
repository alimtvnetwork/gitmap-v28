package message

import "strings"

func applyTitleAffix(msg, prefix, suffix string) string {
	if prefix == "" && suffix == "" {
		return msg
	}
	idx := strings.IndexByte(msg, '\n')
	if idx < 0 {
		return prefix + msg + suffix
	}
	return prefix + msg[:idx] + suffix + msg[idx:]
}

func applyBodyAffix(msg string, prefixPool, suffixPool []string, pick func(int) int) string {
	if len(prefixPool) > 0 {
		msg = pickOne(prefixPool, pick) + "\n" + msg
	}
	if len(suffixPool) > 0 {
		msg = msg + "\n" + pickOne(suffixPool, pick)
	}
	return msg
}

func pickOne(pool []string, pick func(int) int) string {
	if pick == nil {
		return pool[0]
	}
	i := pick(len(pool))
	if i < 0 || i >= len(pool) {
		i = 0
	}
	return pool[i]
}
