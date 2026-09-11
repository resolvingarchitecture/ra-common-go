package config

import "testing"

func TestProperties(t *testing.T) {
	props := ParseProperties("# comment\nfoo=bar\nbaz: qux\n\nempty=\n")
	if props["foo"] != "bar" {
		t.Errorf("foo = %q, want bar", props["foo"])
	}
	if props["baz"] != "qux" {
		t.Errorf("baz = %q, want qux", props["baz"])
	}
	if props["empty"] != "" {
		t.Errorf("empty = %q, want empty string", props["empty"])
	}
}

func TestEnvAndArgs(t *testing.T) {
	env := LoadFromEnv()
	if len(env) == 0 {
		t.Error("LoadFromEnv() should not be empty in a normal process")
	}

	args := LoadFromArgs([]string{"a=1", "b=2", "noeq"}, "=")
	if args["a"] != "1" || args["b"] != "2" {
		t.Errorf("args = %v", args)
	}
	if _, ok := args["noeq"]; ok {
		t.Error("an arg with no delimiter should be dropped")
	}
}

func TestSystemDirs(t *testing.T) {
	if _, err := UserHomeDir(); err != nil {
		t.Errorf("UserHomeDir: %v", err)
	}
	if _, err := UserDataDir(); err != nil {
		t.Errorf("UserDataDir: %v", err)
	}
}
