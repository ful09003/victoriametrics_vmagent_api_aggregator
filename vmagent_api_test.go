package main

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
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
			want: VMAgentAPIResponse{
				Data:   Data{},
				Status: "success",
			},
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
