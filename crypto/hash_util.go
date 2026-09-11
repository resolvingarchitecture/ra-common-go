// Package crypto: hashing, fingerprints, and password hashing. Ports
// ra.common.HashUtil. Uses Go's stdlib crypto/sha1, crypto/sha256,
// crypto/sha512 and crypto/hmac directly - unlike ra-common-cpp, which had to
// implement SHA itself, Go's standard library already covers this. PBKDF2
// isn't in the stdlib (it lives in the golang.org/x/crypto module), but its
// outer loop is a short, well-specified wrapper around crypto/hmac, so it's
// hand-rolled here rather than adding a dependency for ~15 lines.
package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

const (
	pbkdf2Iterations = 1000
	pbkdf2KeyLen     = 64
	saltLen          = 16
)

func Salt() []byte {
	out := make([]byte, saltLen)
	if _, err := rand.Read(out); err != nil {
		panic(err)
	}
	return out
}

func Digest(data []byte, algorithm HashAlgorithm) ([]byte, error) {
	switch algorithm {
	case Sha1:
		h := sha1.Sum(data)
		return h[:], nil
	case Sha256:
		h := sha256.Sum256(data)
		return h[:], nil
	case Sha512:
		h := sha512.Sum512(data)
		return h[:], nil
	default:
		return nil, raerror.CryptoErr("PBKDF2 is not a plain digest")
	}
}

// ToHex returns uppercase hex of data, grouped into blocks of four chars separated by ':'.
func ToHex(data []byte) string {
	hex := fmt.Sprintf("%X", data)
	var groups []string
	for i := 0; i < len(hex); i += 4 {
		end := i + 4
		if end > len(hex) {
			end = len(hex)
		}
		groups = append(groups, hex[i:end])
	}
	return strings.Join(groups, ":")
}

func FromHex(text string) ([]byte, error) {
	clean := strings.ReplaceAll(text, ":", "")
	out := make([]byte, len(clean)/2)
	for i := range out {
		var b int
		if _, err := fmt.Sscanf(clean[i*2:i*2+2], "%02x", &b); err != nil {
			return nil, raerror.DecodeErr("bad hex: " + text)
		}
		out[i] = byte(b)
	}
	return out, nil
}

func GenerateFingerprint(data []byte, algorithm HashAlgorithm) (string, error) {
	d, err := Digest(data, algorithm)
	if err != nil {
		return "", err
	}
	return ToHex(d), nil
}

func pbkdf2HmacSha1(password string, salt []byte, iterations, keyLen int) []byte {
	out := make([]byte, 0, keyLen)
	var blockIndex uint32 = 1
	for len(out) < keyLen {
		blockSuffix := make([]byte, 4)
		binary.BigEndian.PutUint32(blockSuffix, blockIndex)

		mac := hmac.New(sha1.New, []byte(password))
		mac.Write(salt)
		mac.Write(blockSuffix)
		u := mac.Sum(nil)
		t := make([]byte, len(u))
		copy(t, u)

		for i := 1; i < iterations; i++ {
			mac := hmac.New(sha1.New, []byte(password))
			mac.Write(u)
			u = mac.Sum(nil)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		out = append(out, t...)
		blockIndex++
	}
	return out[:keyLen]
}

func GeneratePasswordHash(password string) string {
	return generatePasswordHashWithSalt(password, Salt())
}

func generatePasswordHashWithSalt(password string, salt []byte) string {
	derived := pbkdf2HmacSha1(password, salt, pbkdf2Iterations, pbkdf2KeyLen)
	return fmt.Sprintf("%d_%s_%s", pbkdf2Iterations, base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(derived))
}

func VerifyPasswordHash(password, hashToVerify string) bool {
	parts := strings.Split(hashToVerify, "_")
	if len(parts) != 3 {
		return false
	}
	iterations, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	salt, err1 := base64.StdEncoding.DecodeString(parts[1])
	expected, err2 := base64.StdEncoding.DecodeString(parts[2])
	if err1 != nil || err2 != nil {
		return false
	}
	actual := pbkdf2HmacSha1(password, salt, iterations, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func GenerateHash(content []byte, algorithm HashAlgorithm) (string, error) {
	if algorithm == Pbkdf2HmacSha1 {
		return GeneratePasswordHash(string(content)), nil
	}
	salt := Salt()
	h, err := Digest(append(append([]byte{}, salt...), content...), algorithm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", base64.StdEncoding.EncodeToString(h), base64.StdEncoding.EncodeToString(salt)), nil
}

func VerifyHash(content []byte, hashToVerify string, algorithm HashAlgorithm) (bool, error) {
	if algorithm == Pbkdf2HmacSha1 {
		return VerifyPasswordHash(string(content), hashToVerify), nil
	}
	parts := strings.SplitN(hashToVerify, "_", 2)
	if len(parts) != 2 {
		return false, raerror.InvalidErr("malformed hash")
	}
	expected, err1 := base64.StdEncoding.DecodeString(parts[0])
	salt, err2 := base64.StdEncoding.DecodeString(parts[1])
	if err1 != nil || err2 != nil {
		return false, raerror.InvalidErr("malformed hash")
	}
	actual, err := Digest(append(append([]byte{}, salt...), content...), algorithm)
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}
