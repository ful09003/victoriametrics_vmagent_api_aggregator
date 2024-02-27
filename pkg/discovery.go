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
// <prefix>_VMAGENT_ENDPOINTS env var. Prefix may be empty or set to user preference.
func NewEnvDiscovery(prefix string) *EnvDiscovery {
	return &EnvDiscovery{discoveryEnvVar: prefix + "VMAGENT_ENDPOINTS"}
}

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

type FileDiscovery struct {
	fp string
}

func NewFileDiscovery(fp string) *FileDiscovery {
	return &FileDiscovery{fp: fp}
}

func (f *FileDiscovery) DiscoverEndpoints() ([]string, error) {
	b, err := os.ReadFile(f.fp)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), ","), nil
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
