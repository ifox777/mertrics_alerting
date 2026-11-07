package config_test

import (
	"flag"
	"testing"
	"time"

	"github.com/ifox777/metrics-alerting/internal/config"
	"github.com/stretchr/testify/require"
)

func TestParseServer(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config.ServerConfig
		wantErr bool
	}{
		{
			name: "defaults with nil arguments",
			want: config.ServerConfig{Address: "localhost:8080"},
		},
		{
			name: "defaults with empty arguments",
			args: []string{},
			want: config.ServerConfig{Address: "localhost:8080"},
		},
		{
			name: "custom address",
			args: []string{"-a", "localhost:9090"},
			want: config.ServerConfig{Address: "localhost:9090"},
		},
		{
			name: "address with equals sign",
			args: []string{"-a=127.0.0.1:9091"},
			want: config.ServerConfig{Address: "127.0.0.1:9091"},
		},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "report flag belongs to agent", args: []string{"-r", "3"}, wantErr: true},
		{name: "poll flag belongs to agent", args: []string{"-p", "1"}, wantErr: true},
		{name: "missing address", args: []string{"-a"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ParseServer(tt.args)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseAgent(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config.AgentConfig
		wantErr bool
	}{
		{
			name: "defaults with nil arguments",
			want: config.AgentConfig{
				Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
			},
		},
		{
			name: "defaults with empty arguments",
			args: []string{},
			want: config.AgentConfig{
				Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
			},
		},
		{
			name: "custom address preserves intervals",
			args: []string{"-a", "localhost:9090"},
			want: config.AgentConfig{
				Address: "localhost:9090", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
			},
		},
		{
			name: "custom report interval in seconds",
			args: []string{"-r", "3"},
			want: config.AgentConfig{
				Address: "localhost:8080", ReportInterval: 3 * time.Second, PollInterval: 2 * time.Second,
			},
		},
		{
			name: "custom poll interval in seconds",
			args: []string{"-p", "5"},
			want: config.AgentConfig{
				Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 5 * time.Second,
			},
		},
		{
			name: "all flags",
			args: []string{"-a", "127.0.0.1:9090", "-r", "7", "-p", "3"},
			want: config.AgentConfig{
				Address: "127.0.0.1:9090", ReportInterval: 7 * time.Second, PollInterval: 3 * time.Second,
			},
		},
		{
			name: "all flags with equals signs",
			args: []string{"-p=4", "-r=9", "-a=localhost:9091"},
			want: config.AgentConfig{
				Address: "localhost:9091", ReportInterval: 9 * time.Second, PollInterval: 4 * time.Second,
			},
		},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "nonnumeric report interval", args: []string{"-r", "invalid"}, wantErr: true},
		{name: "nonnumeric poll interval", args: []string{"-p", "invalid"}, wantErr: true},
		{name: "fractional report interval", args: []string{"-r", "1.5"}, wantErr: true},
		{name: "fractional poll interval", args: []string{"-p", "1.5"}, wantErr: true},
		{name: "missing address", args: []string{"-a"}, wantErr: true},
		{name: "missing report interval", args: []string{"-r"}, wantErr: true},
		{name: "missing poll interval", args: []string{"-p"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ParseAgent(tt.args)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseHelp(t *testing.T) {
	parsers := []struct {
		name  string
		parse func([]string) error
	}{
		{
			name: "server",
			parse: func(args []string) error {
				_, err := config.ParseServer(args)
				return err
			},
		},
		{
			name: "agent",
			parse: func(args []string) error {
				_, err := config.ParseAgent(args)
				return err
			},
		},
	}

	for _, parser := range parsers {
		t.Run(parser.name, func(t *testing.T) {
			for _, arg := range []string{"-h", "-help"} {
				t.Run(arg, func(t *testing.T) {
					err := parser.parse([]string{arg})
					require.ErrorIs(t, err, flag.ErrHelp)
				})
			}
		})
	}
}

func TestParseServerRepeatedCalls(t *testing.T) {
	first, err := config.ParseServer([]string{"-a", "localhost:9090"})
	require.NoError(t, err)
	require.Equal(t, config.ServerConfig{Address: "localhost:9090"}, first)

	defaults, err := config.ParseServer(nil)
	require.NoError(t, err)
	require.Equal(t, config.ServerConfig{Address: "localhost:8080"}, defaults)

	_, err = config.ParseServer([]string{"-a", "localhost:9091", "-unknown"})
	require.Error(t, err)

	afterError, err := config.ParseServer(nil)
	require.NoError(t, err)
	require.Equal(t, config.ServerConfig{Address: "localhost:8080"}, afterError)

	require.Equal(t, config.ServerConfig{Address: "localhost:9090"}, first)
}

func TestParseAgentRepeatedCalls(t *testing.T) {
	first, err := config.ParseAgent([]string{"-a", "localhost:9090", "-r", "7", "-p", "3"})
	require.NoError(t, err)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:9090", ReportInterval: 7 * time.Second, PollInterval: 3 * time.Second,
	}, first)

	defaults, err := config.ParseAgent(nil)
	require.NoError(t, err)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
	}, defaults)

	_, err = config.ParseAgent([]string{"-a", "localhost:9091", "-r", "8", "-p", "invalid"})
	require.Error(t, err)

	afterError, err := config.ParseAgent(nil)
	require.NoError(t, err)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
	}, afterError)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:9090", ReportInterval: 7 * time.Second, PollInterval: 3 * time.Second,
	}, first)
}

func TestParsersIndependent(t *testing.T) {
	agent, err := config.ParseAgent([]string{"-a", "localhost:9090", "-r", "7", "-p", "3"})
	require.NoError(t, err)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:9090", ReportInterval: 7 * time.Second, PollInterval: 3 * time.Second,
	}, agent)

	server, err := config.ParseServer([]string{"-a", "localhost:9091"})
	require.NoError(t, err)
	require.Equal(t, config.ServerConfig{Address: "localhost:9091"}, server)

	_, err = config.ParseServer([]string{"-r", "5"})
	require.Error(t, err)

	agentDefaults, err := config.ParseAgent(nil)
	require.NoError(t, err)
	require.Equal(t, config.AgentConfig{
		Address: "localhost:8080", ReportInterval: 10 * time.Second, PollInterval: 2 * time.Second,
	}, agentDefaults)

	serverDefaults, err := config.ParseServer(nil)
	require.NoError(t, err)
	require.Equal(t, config.ServerConfig{Address: "localhost:8080"}, serverDefaults)
}
