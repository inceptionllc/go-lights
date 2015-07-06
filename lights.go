package lights

import "strconv"

// ParseHex parses a hex string.
func ParseHex(val string) (int, error) {
	v, err := strconv.ParseInt(val, 16, 0)
	return int(v), err
}
