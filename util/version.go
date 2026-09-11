package util

// VersionCompare compares two version strings loosely. Returns -1 / 0 / 1.
func VersionCompare(left, right string) int {
	if left == right {
		return 0
	}
	ll, rl := len(left), len(right)
	il, ir := 0, 0

	nextSeparator := func(s string, start int) int {
		i := start
		for i < len(s) && !isSeparator(s[i]) {
			i++
		}
		return i
	}
	parseLong := func(s string, start, end int) int {
		rv := 0
		parsedAny := false
		for i := start; i < end && rv >= 0; i++ {
			c := s[i]
			if c >= '0' && c <= '9' {
				parsedAny = true
				rv = rv*10 + int(c-'0')
			}
		}
		if !parsedAny {
			return -1
		}
		return rv
	}

	for {
		if il >= ll {
			if ir >= rl {
				return 0
			}
			return -1
		}
		if ir >= rl {
			return 1
		}

		lv := -1
		for lv == -1 && il < ll {
			nl := nextSeparator(left, il)
			lv = parseLong(left, il, nl)
			il = nl + 1
		}
		rv := -1
		for rv == -1 && ir < rl {
			nr := nextSeparator(right, ir)
			rv = parseLong(right, ir, nr)
			ir = nr + 1
		}
		if lv < rv {
			return -1
		}
		if lv > rv {
			return 1
		}
	}
}

func isSeparator(c byte) bool {
	return c == '.' || c == '-' || c == '_'
}
