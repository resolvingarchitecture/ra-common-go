// Package content: typed content submitted to the network for dissemination.
// Ports ra.common.content.Content and its subclasses (Text, HTML, JSON,
// Binary, Image, Audio, Video). The Java class hierarchy collapses into one
// Content struct tagged with a Kind.
package content

import (
	"strconv"
	"strings"
	"time"

	"github.com/resolvingarchitecture/ra-common-go/crypto"
	"github.com/resolvingarchitecture/ra-common-go/raerror"
	"github.com/resolvingarchitecture/ra-common-go/util"
)

type Kind string

const (
	Text   Kind = "Text"
	Html   Kind = "Html"
	JSON   Kind = "Json"
	Image  Kind = "Image"
	Audio  Kind = "Audio"
	Video  Kind = "Video"
	Binary Kind = "Binary"
)

func KindForType(contentType string) (Kind, bool) {
	switch {
	case strings.HasPrefix(contentType, "text/plain"):
		return Text, true
	case strings.HasPrefix(contentType, "text/html"):
		return Html, true
	case strings.HasPrefix(contentType, "application/json"):
		return JSON, true
	case strings.HasPrefix(contentType, "image/"):
		return Image, true
	case strings.HasPrefix(contentType, "audio/"):
		return Audio, true
	case strings.HasPrefix(contentType, "video/"):
		return Video, true
	}
	return "", false
}

func KindIsText(k Kind) bool { return k == Text || k == Html || k == JSON }

type Content struct {
	Kind                          Kind                        `json:"type"`
	ContentType                   string                      `json:"content_type"`
	Version                       int                         `json:"version"`
	ID                            *string                     `json:"id,omitempty"`
	Label                         *string                     `json:"label,omitempty"`
	Name                          *string                     `json:"name,omitempty"`
	Location                      *string                     `json:"location,omitempty"`
	Size                          int                         `json:"size"`
	AuthorAlias                   *string                     `json:"author_alias,omitempty"`
	AuthorAddress                 *string                     `json:"author_address,omitempty"`
	Body                          []byte                      `json:"body,omitempty"`
	BodyEncoding                  *string                     `json:"body_encoding,omitempty"`
	BodyBase64Encoded             bool                        `json:"body_base64_encoded"`
	CreatedAt                     *int64                      `json:"created_at,omitempty"`
	Hash                          *crypto.Hash                `json:"hash,omitempty"`
	HashAlgorithm                 crypto.HashAlgorithm        `json:"hash_algorithm"`
	Fingerprint                   *crypto.Hash                `json:"fingerprint,omitempty"`
	FingerprintAlgorithm          crypto.HashAlgorithm        `json:"fingerprint_algorithm"`
	Children                      []Content                   `json:"children,omitempty"`
	Encrypted                     bool                        `json:"encrypted"`
	EncryptionAlgorithm           *crypto.EncryptionAlgorithm `json:"encryption_algorithm,omitempty"`
	EncryptionPassphrase          *string                     `json:"encryption_passphrase,omitempty"`
	EncryptionPassphraseEncrypted bool                        `json:"encryption_passphrase_encrypted"`
	EncryptionPassphraseAlgorithm *crypto.EncryptionAlgorithm `json:"encryption_passphrase_algorithm,omitempty"`
	Base64EncodedIv               *string                     `json:"base64_encoded_iv,omitempty"`
	Keywords                      []string                    `json:"keywords,omitempty"`
	Readable                      bool                        `json:"readable"`
	Writeable                     bool                        `json:"writeable"`
}

func New(kind Kind, contentType string) Content {
	return Content{
		Kind:                 kind,
		ContentType:          contentType,
		HashAlgorithm:        crypto.Sha256,
		FingerprintAlgorithm: crypto.Sha1,
	}
}

type BuildOptions struct {
	Label               *string
	Name                *string
	GenerateHash        bool
	GenerateFingerprint bool
}

func Build(body []byte, contentType string, opts BuildOptions) (*Content, error) {
	kind, ok := KindForType(contentType)
	if !ok {
		return nil, raerror.InvalidErr("unsupported content type: " + contentType)
	}
	c := New(kind, contentType)
	c.Label = opts.Label
	c.Name = opts.Name
	if idx := strings.Index(contentType, "charset:"); idx >= 0 {
		enc := contentType[idx+len("charset:"):]
		c.BodyEncoding = &enc
	}
	if err := c.SetBody(body, opts.GenerateHash, opts.GenerateFingerprint); err != nil {
		return nil, err
	}
	createdAt := time.Now().UnixMilli()
	c.CreatedAt = &createdAt
	id := util.RandomAlphanumeric(32)
	c.ID = &id
	return &c, nil
}

func (c *Content) SetBody(body []byte, generateHash, generateFingerprint bool) error {
	c.Size = len(body)
	if generateHash {
		h, err := crypto.GenerateHash(body, c.HashAlgorithm)
		if err != nil {
			return err
		}
		c.Hash = &crypto.Hash{HashValue: h, Algorithm: c.HashAlgorithm}
	}
	if generateFingerprint && c.Hash != nil {
		fp, err := crypto.GenerateFingerprint([]byte(c.Hash.HashValue), c.FingerprintAlgorithm)
		if err != nil {
			return err
		}
		c.Fingerprint = &crypto.Hash{HashValue: fp, Algorithm: c.FingerprintAlgorithm}
	}
	c.Body = body
	c.Version++
	return nil
}

func (c *Content) MetaOnly() bool { return c.Body == nil }

func (c *Content) AddKeyword(keyword string) { c.Keywords = append(c.Keywords, keyword) }
func (c *Content) AddChild(child Content)    { c.Children = append(c.Children, child) }

func (c *Content) MagnetLink() *string {
	var parts []string
	if c.Body != nil {
		parts = append(parts, "xl="+strconv.Itoa(len(c.Body)))
	}
	if c.Hash != nil {
		parts = append(parts, "xt=urn:"+strings.ToLower(c.Hash.Algorithm.JcaName())+":"+c.Hash.HashValue)
	}
	if len(c.Keywords) > 0 {
		parts = append(parts, "kt="+strings.Join(c.Keywords, "+"))
	}
	if len(parts) == 0 {
		return nil
	}
	link := "magnet:?" + strings.Join(parts, "&")
	return &link
}
