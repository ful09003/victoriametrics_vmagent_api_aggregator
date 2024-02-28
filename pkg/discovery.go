package pkg

import (
	"fmt"
	"os"
	"strings"
)

type VMAgentDiscoverer interface {
	DiscoverEndpoints() ([]string, error)
}

// EnvDiscovery returns vmagent endpoints defined in the environment.
type EnvDiscovery struct {
	discoveryEnvVar string
}

// NewEnvDiscovery will return a new EnvDiscovery instance which will look up vmagent endpoints from the
// provided env var.
func NewEnvDiscovery(envvar string) *EnvDiscovery {
	return &EnvDiscovery{discoveryEnvVar: envvar}
}

// DiscoverEndpoints returns a string slice for endpoints from the EnvDiscovery configured env var
// Because env vars cannot change post-process-start, an empty env var will return an error
func (e *EnvDiscovery) DiscoverEndpoints() ([]string, error) {
	eV := gatherEnvVars(os.Environ(), []string{e.discoveryEnvVar})
	val, ok := eV[e.discoveryEnvVar]
	if !ok {
		return nil, fmt.Errorf("no endpoints discovered. double check env var %s", e.discoveryEnvVar)
	}
	if val == "" {
		return nil, fmt.Errorf("no endpoints defined. double check env var %s", e.discoveryEnvVar)
	}
	return strings.Split(val, ","), nil
}

// FileDiscovery returns vmagent endpoints based on file contents
type FileDiscovery struct {
	fp string
}

func NewFileDiscovery(fp string) *FileDiscovery {
	return &FileDiscovery{fp: fp}
}

// DiscoverEndpoints returns a slice of strings hopefully containing vmagent endpoints, from the FileDiscovery's configured
// file.
func (f *FileDiscovery) DiscoverEndpoints() ([]string, error) {
	b, err := os.ReadFile(f.fp)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(b)), ","), nil
}

func gatherEnvVars(inEnv []string, k []string) map[string]string {
	ret := make(map[string]string)
	for _, desiredKey := range k {
		k, v := gatherVarVal(inEnv, desiredKey)
		ret[k] = v
	}

	return ret
}

func gatherVarVal(inEnv []string, k string) (string, string) {
	for _, envVarRaw := range inEnv {
		before, after, found := strings.Cut(envVarRaw, "=")
		if found && before == k {
			return k, after
		}
	}
	return k, ""
}
