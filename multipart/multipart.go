// Package multipart ports ra.common.file.Multipart - a multipart/form-data
// body builder. Like the Java version (whose HTTP transport was commented
// out) this only accumulates the body string; sending it is the caller's concern.
package multipart

import (
	"fmt"
	"strings"

	"github.com/resolvingarchitecture/ra-common-go/util"
)

const lineFeed = "\r\n"

type Multipart struct {
	Boundary string  `json:"boundary"`
	Charset  *string `json:"charset,omitempty"`
	body     strings.Builder
}

func New(charset *string, boundary *string) *Multipart {
	b := ""
	if boundary != nil {
		b = *boundary
	} else {
		b = util.RandomAlphanumeric(32)
	}
	return &Multipart{Boundary: b, Charset: charset}
}

func (m *Multipart) AddFormField(name, value string) {
	charset := "UTF-8"
	if m.Charset != nil {
		charset = *m.Charset
	}
	fmt.Fprintf(&m.body, "--%s%s", m.Boundary, lineFeed)
	fmt.Fprintf(&m.body, "Content-Disposition: form-data; name=\"%s\"%s", name, lineFeed)
	fmt.Fprintf(&m.body, "Content-Type: text/plain; charset=%s%s%s", charset, lineFeed, lineFeed)
	fmt.Fprintf(&m.body, "%s%s", value, lineFeed)
}

func (m *Multipart) AddFilePart(fieldName string, fileName *string) {
	fmt.Fprintf(&m.body, "--%s%s", m.Boundary, lineFeed)
	if fileName != nil {
		fmt.Fprintf(&m.body, "Content-Disposition: file; filename=\"%s\"%s", *fileName, lineFeed)
	} else {
		fmt.Fprintf(&m.body, "Content-Disposition: file; name=\"%s\";%s", fieldName, lineFeed)
	}
	fmt.Fprintf(&m.body, "Content-Type: application/octet-stream%s", lineFeed)
	fmt.Fprintf(&m.body, "Content-Transfer-Encoding: binary%s%s", lineFeed, lineFeed)
}

func (m *Multipart) AddHeaderField(name, value string) {
	fmt.Fprintf(&m.body, "%s: %s%s", name, value, lineFeed)
}

func (m *Multipart) AppendRaw(text string) { m.body.WriteString(text) }

func (m *Multipart) Body() string { return m.body.String() }

func (m *Multipart) Finish() string {
	return m.body.String() + fmt.Sprintf("--%s--%s", m.Boundary, lineFeed)
}
