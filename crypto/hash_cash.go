// Package crypto: Hashcash proof-of-work tokens. Ports ra.common.HashCash.
package crypto

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/bits"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

const hashBits = 160

type HashCash struct {
	token      string
	value      int
	resource   string
	date       time.Time // calendar date only (Y/M/D); time-of-day and location are not meaningful here
	version    int
	extensions map[string][]string
}

func (h *HashCash) Token() string                   { return h.token }
func (h *HashCash) Value() int                      { return h.value }
func (h *HashCash) Resource() string                { return h.resource }
func (h *HashCash) Date() time.Time                 { return h.date }
func (h *HashCash) Version() int                    { return h.version }
func (h *HashCash) Extensions() map[string][]string { return h.extensions }

func MintHashCash(resource string, bitsWanted int) (*HashCash, error) {
	return MintHashCashWith(resource, map[string][]string{}, time.Now().UTC(), bitsWanted, 1)
}

func MintHashCashWith(resource string, extensions map[string][]string, date time.Time, bitsWanted, version int) (*HashCash, error) {
	if version > 1 {
		return nil, raerror.InvalidErr("only hashcash versions 0 and 1 are supported")
	}
	if bitsWanted > hashBits {
		return nil, raerror.InvalidErr("value must be between 0 and 160")
	}
	if strings.Contains(resource, ":") {
		return nil, raerror.InvalidErr("resource may not contain a colon")
	}
	extStr, err := serializeExtensions(extensions)
	if err != nil {
		return nil, err
	}
	dateStr := fmtYyMmDd(date)
	var prefix string
	if version == 0 {
		prefix = fmt.Sprintf("0:%s:%s:%s:", dateStr, resource, extStr)
	} else {
		prefix = fmt.Sprintf("1:%d:%s:%s:%s:", bitsWanted, dateStr, resource, extStr)
	}
	token := generateHashCashToken(prefix, bitsWanted)
	value := bitsWanted
	if version == 0 {
		value = sha1Bits(token)
	}
	return &HashCash{token: token, value: value, resource: resource, date: date, version: version, extensions: extensions}, nil
}

func generateHashCashToken(prefix string, bitsWanted int) string {
	rndBytes := make([]byte, 8)
	_, _ = rand.Read(rndBytes)
	rnd := fmt.Sprintf("%x", rndBytes)

	counterBytes := make([]byte, 8)
	_, _ = rand.Read(counterBytes)
	counter := binary.BigEndian.Uint64(counterBytes)

	stem := prefix + rnd + ":"
	for {
		counter++
		candidate := fmt.Sprintf("%s%x", stem, counter)
		if sha1Bits(candidate) >= bitsWanted {
			return candidate
		}
	}
}

func ParseHashCash(token string) (*HashCash, error) {
	parts := strings.Split(token, ":")
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, raerror.InvalidErr("bad hashcash version")
	}
	var expected int
	switch version {
	case 0:
		expected = 6
	case 1:
		expected = 7
	default:
		return nil, raerror.InvalidErr("only hashcash versions 0 and 1 are supported")
	}
	if len(parts) != expected {
		return nil, raerror.InvalidErr("improperly formed hashcash")
	}

	idx := 1
	claimedBits := 0
	if version == 1 {
		claimedBits, err = strconv.Atoi(parts[idx])
		if err != nil {
			return nil, raerror.InvalidErr("bad hashcash bits")
		}
		idx++
	}
	date, err := parseYyMmDd(parts[idx])
	if err != nil {
		return nil, err
	}
	idx++
	resource := parts[idx]
	idx++
	extensions := deserializeExtensions(parts[idx])

	actual := sha1Bits(token)
	value := actual
	if version != 0 && claimedBits < actual {
		value = claimedBits
	}
	return &HashCash{token: token, value: value, resource: resource, date: date, version: version, extensions: extensions}, nil
}

func (h *HashCash) ComputedBits() int { return sha1Bits(h.token) }

func (h *HashCash) IsValidFor(resource string, minBits int) bool {
	return h.resource == resource && h.ComputedBits() >= minBits
}

func (h *HashCash) String() string { return h.token }

func sha1Bits(token string) int {
	sum := sha1.Sum([]byte(token))
	return leadingZeroBits(sum[:])
}

// LeadingZeroBits is exposed for testing, mirroring leadingZeroBitsForTest in ra-common-ts.
func LeadingZeroBits(data []byte) int { return leadingZeroBits(data) }

func leadingZeroBits(data []byte) int {
	total := 0
	for _, b := range data {
		if b == 0 {
			total += 8
		} else {
			total += bits.LeadingZeros8(b)
			break
		}
	}
	return total
}

func serializeExtensions(ext map[string][]string) (string, error) {
	if len(ext) == 0 {
		return "", nil
	}
	keys := make([]string, 0, len(ext))
	for k := range ext {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, key := range keys {
		if strings.ContainsAny(key, ":;=") {
			return "", raerror.InvalidErr("illegal char in extension key: " + key)
		}
		values := ext[key]
		if len(values) == 0 {
			parts = append(parts, key)
			continue
		}
		for _, v := range values {
			if strings.ContainsAny(v, ":;,") {
				return "", raerror.InvalidErr("illegal char in extension value: " + v)
			}
		}
		parts = append(parts, key+"="+strings.Join(values, ","))
	}
	return strings.Join(parts, ";"), nil
}

func deserializeExtensions(text string) map[string][]string {
	out := map[string][]string{}
	if text == "" {
		return out
	}
	for _, item := range strings.Split(text, ";") {
		if eq := strings.Index(item, "="); eq == -1 {
			out[item] = []string{}
		} else {
			out[item[:eq]] = strings.Split(item[eq+1:], ",")
		}
	}
	return out
}

func fmtYyMmDd(d time.Time) string {
	return fmt.Sprintf("%02d%02d%02d", d.Year()%100, int(d.Month()), d.Day())
}

func parseYyMmDd(s string) (time.Time, error) {
	if len(s) != 6 {
		return time.Time{}, raerror.InvalidErr("bad hashcash date: " + s)
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return time.Time{}, raerror.InvalidErr("bad hashcash date: " + s)
		}
	}
	yy, _ := strconv.Atoi(s[0:2])
	mm, _ := strconv.Atoi(s[2:4])
	dd, _ := strconv.Atoi(s[4:6])
	date := time.Date(2000+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	if date.Year() != 2000+yy || int(date.Month()) != mm || date.Day() != dd {
		return time.Time{}, raerror.InvalidErr("bad hashcash date: " + s)
	}
	return date, nil
}
