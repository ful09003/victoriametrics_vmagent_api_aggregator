package pkg

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_gatherVarVal(t *testing.T) {
	type args struct {
		inEnv []string
		k     string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "happy path",
			args: args{
				inEnv: []string{"SHELL=/bin/bash"},
				k:     "SHELL",
			},
			want:  "SHELL",
			want1: "/bin/bash",
		},
		{
			name: "happy path not found",
			args: args{
				inEnv: []string{"SHELL=/bin/bash"},
				k:     "USER",
			},
			want:  "USER",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := gatherVarVal(tt.args.inEnv, tt.args.k)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.want1)
		})
	}
}

func Test_GatherEnvVar(t *testing.T) {
	type args struct {
		inEnv []string
		k     []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "happy path",
			args: args{
				inEnv: []string{"SHELL=/bin/bash", "USER=me"},
				k:     []string{"SHELL", "USER"},
			},
			want: map[string]string{"SHELL": "/bin/bash", "USER": "me"},
		},
		{
			name: "happy path nothing found",
			args: args{
				inEnv: []string{"SHELL=/bin/bash", "USER=me"},
				k:     []string{"TERM"},
			},
			want: map[string]string{"TERM": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.DeepEqual(t, gatherEnvVars(tt.args.inEnv, tt.args.k), tt.want)
		})
	}
}

func TestEnvDiscovery_DiscoverEndpoints(t *testing.T) {
	type fields struct {
		envVarsToSet map[string]string
		envVar       string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "happy path with unset prefix",
			fields: fields{
				envVarsToSet: map[string]string{"VMAGENT_ENDPOINTS": "http://1.2.3.4:1234"},
				envVar:       "VMAGENT_ENDPOINTS",
			},
			want:    []string{"http://1.2.3.4:1234"},
			wantErr: false,
		},
		{
			name: "happy path errors when env not found",
			fields: fields{
				envVarsToSet: nil,
				envVar:       "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path errors when env var empty",
			fields: fields{
				envVarsToSet: map[string]string{"VMAGENT_ENDPOINTS": ""},
				envVar:       "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.fields.envVarsToSet {
				t.Setenv(k, v)
			}
			e := &EnvDiscovery{
				discoveryEnvVar: tt.fields.envVar,
			}
			got, err := e.DiscoverEndpoints()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscoverEndpoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.DeepEqual(t, got, tt.want)
		})
	}
}
