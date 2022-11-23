package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chroju/parade/ssmctl"
)

const (
	usageDelHelp = "  -h, --help    help for del\n\n"
)

func Test_delCommand(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		{
			name:          "one arg",
			command:       "/service1/dev/key1",
			wantOutWriter: "",
			wantErrWriter: "Delete `/service1/dev/key1` (value: dev_value1) ? (Y/n)\n",
			wantErr:       false,
		},
		{
			name:          "one arg with force option",
			command:       "/service1/dev/key1 --force",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "two args",
			command:       "dev prod",
			wantOutWriter: "",
			wantErrWriter: "accepts 1 arg(s), received 2",
			wantErr:       true,
		},
		{
			name:          "no args",
			command:       "",
			wantOutWriter: "",
			wantErrWriter: "accepts 1 arg(s), received 0",
			wantErr:       true,
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

			cmd := newDelCommand(o)
			if tt.command != "" {
				cmd.SetArgs(strings.Split(tt.command, " "))
			}

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err.Error() != tt.wantErrWriter {
					t.Errorf("get() = %v, want %v", err.Error(), tt.wantErrWriter)
				}
			} else {
				if gotOutWriter := outWriter.String(); gotOutWriter != tt.wantOutWriter {
					t.Errorf("get() = %v, want %v", gotOutWriter, tt.wantOutWriter)
				}
			}
		})
	}
}
