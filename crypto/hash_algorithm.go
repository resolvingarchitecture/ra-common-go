package crypto

import "github.com/resolvingarchitecture/ra-common-go/raerror"

type HashAlgorithm string

const (
	Sha1           HashAlgorithm = "Sha1"
	Sha256         HashAlgorithm = "Sha256"
	Sha512         HashAlgorithm = "Sha512"
	Pbkdf2HmacSha1 HashAlgorithm = "Pbkdf2HmacSha1"
)

func (a HashAlgorithm) JcaName() string {
	switch a {
	case Sha1:
		return "SHA-1"
	case Sha256:
		return "SHA-256"
	case Sha512:
		return "SHA-512"
	case Pbkdf2HmacSha1:
		return "PBKDF2WithHmacSHA1"
	}
	return "SHA-256"
}

func ParseHashAlgorithm(text string) (HashAlgorithm, error) {
	switch text {
	case "SHA-1", "SHA1", "Sha1":
		return Sha1, nil
	case "SHA-256", "SHA256", "Sha256":
		return Sha256, nil
	case "SHA-512", "SHA512", "Sha512":
		return Sha512, nil
	case "PBKDF2WithHmacSHA1", "Pbkdf2HmacSha1":
		return Pbkdf2HmacSha1, nil
	}
	return "", raerror.InvalidErr("unknown hash algorithm: " + text)
}
