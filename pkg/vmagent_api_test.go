package pkg

import (
	"net/http"
	"sync"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_fetchVMAgentTargets(t *testing.T) {
	type args struct {
		c *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    VMAgentAPIResponse
		wantErr error
	}{
		{
			name: "happy path",
			args: args{
				c: http.DefaultClient,
			},
			want:    happyResponse(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, req := genTestHappyPathServerReq(t)
			defer srv.Close()

			got, err := fetchVMAgentTargets(tt.args.c, req)
			assert.Equal(t, err, tt.wantErr)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}

func TestVMAgentAPICollector_Collect(t *testing.T) {
	tests := []struct {
		name    string
		client  *http.Client
		want    VMAgentAPIResponse
		wantErr error
	}{
		{
			name:    "happy",
			client:  http.DefaultClient,
			want:    happyResponse(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, req := genTestHappyPathServerReq(t)
			defer srv.Close()

			v := &VMAgentAPICollector{
				origEndpoint: req.URL.String(),
				c:            tt.client,
			}
			got, err := v.Collect()
			assert.Equal(t, err, tt.wantErr)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}

func TestVMAgentAPICollection_CollectAll(t *testing.T) {
	srv, req := genTestHappyPathServerReq(t)
	defer srv.Close()

	type fields struct {
		m    *sync.Mutex
		c    map[string]*VMAgentAPICollector
		data map[string]VMAgentAPIResponse
	}
	tests := []struct {
		name       string
		fields     fields
		wantErrLen int
	}{
		{
			name: "happy path",
			fields: fields{
				m: &sync.Mutex{},
				c: map[string]*VMAgentAPICollector{
					"local": &VMAgentAPICollector{
						origEndpoint: req.URL.String(),
						c:            http.DefaultClient,
					},
				},
				data: map[string]VMAgentAPIResponse{},
			},
			wantErrLen: 0,
		},
		{
			name: "happy path with a broken vmagent",
			fields: fields{
				m: &sync.Mutex{},
				c: map[string]*VMAgentAPICollector{
					"happy": &VMAgentAPICollector{
						origEndpoint: req.URL.String(),
						c:            http.DefaultClient,
					},
					"sad": &VMAgentAPICollector{
						origEndpoint: "busted",
						c:            http.DefaultClient,
					},
				},
				data: map[string]VMAgentAPIResponse{},
			},
			wantErrLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VMAgentAPICollection{
				m:    tt.fields.m,
				c:    tt.fields.c,
				data: tt.fields.data,
			}
			assert.Equal(t, len(v.CollectAll()), tt.wantErrLen)
		})
	}
}

func TestNewVMAgentCollection(t *testing.T) {
	type args struct {
		disco VMAgentDiscoverer
	}
	tests := []struct {
		name    string
		args    args
		want    *VMAgentAPICollection
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				disco: &StaticMemDiscovery{e: []string{"http://localhost:1234"}},
			},
			want: &VMAgentAPICollection{
				m: &sync.Mutex{},
				c: map[string]*VMAgentAPICollector{
					"http://localhost:1234": {},
				},
				data: map[string]VMAgentAPIResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVMAgentCollection(tt.args.disco)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVMAgentCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(VMAgentAPICollection{}))

		})
	}
}

func TestVMAgentAPICollection_Reconcile(t *testing.T) {
	collection := &VMAgentAPICollection{
		m: &sync.Mutex{},
		c: map[string]*VMAgentAPICollector{
			"http://localhost:1234": &VMAgentAPICollector{},
			"http://localhost:2345": &VMAgentAPICollector{},
		},
		data: map[string]VMAgentAPIResponse{},
	}
	err := collection.Reconcile([]string{"http://localhost:1234", "http://localhost:4567"})
	assert.NilError(t, err)
	assert.DeepEqual(t, collection.c, map[string]*VMAgentAPICollector{
		"http://localhost:1234": &VMAgentAPICollector{},
		"http://localhost:4567": &VMAgentAPICollector{},
	}, cmpopts.IgnoreUnexported(VMAgentAPICollection{}), cmpopts.IgnoreUnexported(VMAgentAPICollector{}))
}
