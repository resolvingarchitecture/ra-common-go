package content

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestContentBuildAndMagnet(t *testing.T) {
	name := "greeting"
	c, err := Build([]byte("hello"), "text/plain", BuildOptions{Name: &name, GenerateHash: true, GenerateFingerprint: true})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if c.Kind != Text {
		t.Errorf("Kind = %q, want Text", c.Kind)
	}
	if c.Size != 5 {
		t.Errorf("Size = %d, want 5", c.Size)
	}
	if c.ID == nil || c.Hash == nil || c.Fingerprint == nil {
		t.Error("ID, Hash and Fingerprint should all be set")
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var back Content
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back.Kind != Text {
		t.Errorf("round-tripped Kind = %q, want Text", back.Kind)
	}
	if !bytes.Equal(back.Body, []byte("hello")) {
		t.Errorf("round-tripped Body = %q, want hello", back.Body)
	}

	b := New(Binary, "application/octet-stream")
	if err := b.SetBody([]byte{1, 2, 3, 4}, true, false); err != nil {
		t.Fatalf("SetBody: %v", err)
	}
	b.AddKeyword("alpha")
	b.AddKeyword("beta")
	link := b.MagnetLink()
	if link == nil {
		t.Fatal("MagnetLink should not be nil")
	}
	if !strings.HasPrefix(*link, "magnet:?xl=4") {
		t.Errorf("MagnetLink = %q, want prefix magnet:?xl=4", *link)
	}
	if !strings.Contains(*link, "kt=alpha+beta") {
		t.Errorf("MagnetLink = %q, want to contain kt=alpha+beta", *link)
	}
}
