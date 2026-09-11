package crypto

type EncryptionAlgorithm string

const (
	Cast5  EncryptionAlgorithm = "Cast5"
	Aes256 EncryptionAlgorithm = "Aes256"
	Aes512 EncryptionAlgorithm = "Aes512"
)

func (a EncryptionAlgorithm) Name() string {
	switch a {
	case Cast5:
		return "CAST-5"
	case Aes256:
		return "AES-256"
	case Aes512:
		return "AES-512"
	}
	return "AES-256"
}
