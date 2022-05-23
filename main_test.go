package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newDemoCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
		wantErr    bool
		errMsg     string
	}{
		{
			name: "empty",
			args: []string{},
			wantStdout: heredoc(`
				LONG TEXT

				Usage:
				  demo [command]

				Available Commands:
				  completion  Generate the autocompletion script for the specified shell
				  hello       SHORT TEXT
				  help        Help about any command

				Flags:
				  -h, --help   help for demo

				Use "demo [command] --help" for more information about a command.
			`),
			wantStderr: "",
			wantErr:    false,
		},
		{
			name: "help flag",
			args: []string{"--help"},
			wantStdout: heredoc(`
				LONG TEXT

				Usage:
				  demo [command]

				Available Commands:
				  completion  Generate the autocompletion script for the specified shell
				  hello       SHORT TEXT
				  help        Help about any command

				Flags:
				  -h, --help   help for demo

				Use "demo [command] --help" for more information about a command.
			`),
			wantStderr: "",
			wantErr:    false,
		},
		{
			name:       "subcommand",
			args:       []string{"hello"},
			wantStdout: "hello world\n",
			wantStderr: "",
			wantErr:    false,
		},
		{
			name: "subcommand help",
			args: []string{"help", "hello"},
			wantStdout: heredoc(`
				LONG TEXT
				
				Usage:
				  demo hello [flags]
				
				Flags:
				  -h, --help      help for hello
				      --num int   number
			`),
			wantStderr: "",
			wantErr:    false,
		},
		{
			name:       "help unknown command",
			args:       []string{"help", "nonexist"},
			wantStdout: "",
			wantStderr: heredoc(`
				Error: unknown command "nonexist" for "demo"
				Usage:
				  demo help [command] [flags]
				
				Flags:
				  -h, --help   help for help
				
			`),
			wantErr: true,
			errMsg:  "unknown command \"nonexist\" for \"demo\"",
		},
		{
			name:       "unknown command",
			args:       []string{"nonexist"},
			wantStdout: "",
			wantStderr: heredoc(`
				Error: unknown command "nonexist" for "demo"
				Run 'demo --help' for usage.
			`),
			wantErr: true,
			errMsg:  `unknown command "nonexist" for "demo"`,
		},
		{
			name:       "deprecated flag",
			args:       []string{"hello", "--flag"},
			wantStdout: "hello world\n",
			wantStderr: "Flag --flag has been deprecated, please don't use it\n",
			wantErr:    false,
		},
		{
			name:       "invalid argument",
			args:       []string{"hello", "nonexist"},
			wantStdout: "",
			wantStderr: heredoc(`
				Error: unknown command "nonexist" for "demo hello"
				Usage:
				  demo hello [flags]
				
				Flags:
				  -h, --help      help for hello
				      --num int   number
				
			`),
			wantErr: true,
			errMsg:  `unknown command "nonexist" for "demo hello"`,
		},
		{
			name:       "invalid flag value",
			args:       []string{"hello", "--num=true"},
			wantStdout: "",
			wantStderr: heredoc(`
				Error: invalid argument "true" for "--num" flag: strconv.ParseInt: parsing "true": invalid syntax
				Usage:
				  demo hello [flags]
				
				Flags:
				  -h, --help      help for hello
				      --num int   number
				
			`),
			wantErr: true,
			errMsg:  `invalid argument "true" for "--num" flag: strconv.ParseInt: parsing "true": invalid syntax`,
		},
		{
			name:       "deprecated command",
			args:       []string{"hullo"},
			wantStdout: "hullo wurld\n",
			wantStderr: "Command \"hullo\" is deprecated, use hello instead\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newDemoCommand()
			c.SetArgs(tt.args)
			stdout := bytes.Buffer{}
			c.SetOut(&stdout)
			stderr := bytes.Buffer{}
			c.SetErr(&stderr)
			_, err := c.ExecuteC()
			if tt.wantErr {
				if err == nil {
					t.Error("ExecuteC() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("ExecuteC() expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("ExecuteC() unexpected error: %v", err)
			}
			assert.Equal(t, tt.wantStdout, stdout.String(), "stdout did not match")
			assert.Equal(t, tt.wantStderr, stderr.String(), "stderr did not match")
		})
	}
}

var tabRE = regexp.MustCompile(`(?m)^\t+`)

// heredoc strips leading tabs from a string and replaces '' with a literal backtick
func heredoc(s string) string {
	s = tabRE.ReplaceAllLiteralString(s, "")
	s = strings.TrimPrefix(s, "\n")
	s = strings.ReplaceAll(s, "''", "`")
	return s
}
