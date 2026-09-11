// Package config: configuration loading and cross-platform directory
// resolution. Ports ra.common.Config and ra.common.SystemSettings. Unlike
// ra-common-cpp (POSIX-only), os.UserHomeDir() and os.Environ() are already
// cross-platform in Go's standard library - no separate Windows backend needed.
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/resolvingarchitecture/ra-common-go/raerror"
)

func LoadFromArgs(args []string, delimiter string) map[string]string {
	out := map[string]string{}
	for _, arg := range args {
		idx := strings.Index(arg, delimiter)
		if idx > 0 {
			out[arg[:idx]] = arg[idx+len(delimiter):]
		}
	}
	return out
}

func LoadFromEnv() map[string]string {
	out := map[string]string{}
	for _, kv := range os.Environ() {
		if idx := strings.Index(kv, "="); idx >= 0 {
			out[kv[:idx]] = kv[idx+1:]
		}
	}
	return out
}

func ParseProperties(text string) map[string]string {
	out := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		for _, sep := range []string{"=", ":"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				out[strings.TrimSpace(line[:idx])] = strings.TrimSpace(line[idx+1:])
				break
			}
		}
	}
	return out
}

func LoadFromFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, raerror.New(raerror.IO, err.Error())
	}
	return ParseProperties(string(data)), nil
}

func LoadAll(clientProps map[string]string, configPath *string) (map[string]string, error) {
	config := LoadFromEnv()
	if configPath != nil {
		fromFile, err := LoadFromFile(*configPath)
		if err != nil {
			return nil, err
		}
		for k, v := range fromFile {
			config[k] = v
		}
	}
	for k, v := range clientProps {
		config[k] = v
	}
	return config, nil
}

// UserHomeDir, UserConfigDir, UserCacheDir wrap Go's already-cross-platform
// os.UserHomeDir/os.UserConfigDir/os.UserCacheDir.
func UserHomeDir() (string, error) { return os.UserHomeDir() }

func UserDataDir() (string, error) {
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}

func UserConfigDir() (string, error) { return os.UserConfigDir() }

func UserCacheDir() (string, error) { return os.UserCacheDir() }

func AppDir(base, group, app string, create bool) (string, error) {
	dir := filepath.Join(base, group, app)
	if create {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", raerror.New(raerror.FileCreationFailed, dir+": "+err.Error())
		}
	}
	return dir, nil
}

func UserAppDataDir(group, app string, create bool) (string, error) {
	base, err := UserDataDir()
	if err != nil {
		return "", raerror.New(raerror.FileCreationFailed, "no user data dir")
	}
	return AppDir(base, group, app, create)
}

func UserAppConfigDir(group, app string, create bool) (string, error) {
	base, err := UserConfigDir()
	if err != nil {
		return "", raerror.New(raerror.FileCreationFailed, "no user config dir")
	}
	return AppDir(base, group, app, create)
}

func UserAppCacheDir(group, app string, create bool) (string, error) {
	base, err := UserCacheDir()
	if err != nil {
		return "", raerror.New(raerror.FileCreationFailed, "no user cache dir")
	}
	return AppDir(base, group, app, create)
}
