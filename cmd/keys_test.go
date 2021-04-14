package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chroju/parade/ssmctl"
)

func Test_keysCommand(t *testing.T) {
	type args struct {
		args       []string
		ssmManager ssmctl.SSMManager
	}
	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		{
			name:          "one arg for Equals",
			command:       "/service1/dev/key1",
			wantOutWriter: "/service1/dev/key1  Type: String\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for BeginsWith",
			command:       "/service1/*",
			wantOutWriter: "/service1/dev/key1   Type: String\n/service1/dev/key2   Type: SecureString\n/service1/prod/key1  Type: String\n/service1/prod/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for Contains",
			command:       "*dev*",
			wantOutWriter: "/service1/dev/key1  Type: String\n/service1/dev/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "one arg for no match",
			command:       "no_match",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "no args",
			command:       "",
			wantOutWriter: "/service1/dev/key1   Type: String\n/service1/dev/key2   Type: SecureString\n/service1/prod/key1  Type: String\n/service1/prod/key2  Type: SecureString\n",
			wantErrWriter: "",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssmManager := ssmctl.NewMockSSMManager()
			outWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			o := &GlobalOption{
				SSMManager: ssmManager,
				Out:        outWriter,
				ErrOut:     errWriter,
			}

			cmd := newKeysCommand(o)
			if tt.command != "" {
				cmd.SetArgs(strings.Split(tt.command, " "))
			}

			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutWriter := outWriter.String(); gotOutWriter != tt.wantOutWriter {
				t.Errorf("keys() = %v, want %v", gotOutWriter, tt.wantOutWriter)
			}
			if gotErrWriter := errWriter.String(); gotErrWriter != tt.wantErrWriter {
				t.Errorf("keys() = %v, want %v", gotErrWriter, tt.wantErrWriter)
			}
		})
	}
}
