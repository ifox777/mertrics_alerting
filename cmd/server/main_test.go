package main

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	serverHelperEnv  = "GO_ADVANCED_SERVER_TEST_PROCESS"
	serverAddressEnv = "GO_ADVANCED_SERVER_TEST_ADDRESS"
	processTimeout   = 8 * time.Second
)

func TestServerProcessHelper(t *testing.T) {
	if os.Getenv(serverHelperEnv) != "1" {
		return
	}

	os.Args = []string{"server", "-a", os.Getenv(serverAddressEnv)}
	main()
	os.Exit(0)
}

func TestServerStartupError(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })

	process := startServerProcess(t, occupied.Addr().String())
	process.wait(t)

	var exitError *exec.ExitError
	require.ErrorAs(t, process.waitErr, &exitError)
	require.Equal(t, 1, exitError.ExitCode())

	output := process.output.String()
	require.Contains(t, output, occupied.Addr().String())
	require.Contains(t, output, "bind")
	require.NotContains(t, output, "panic:")
	require.NotContains(t, output, "\ngoroutine ")
}

func TestServerSignals(t *testing.T) {
	tests := []struct {
		name   string
		signal os.Signal
	}{
		{name: "SIGINT", signal: os.Interrupt},
		{name: "SIGTERM", signal: syscall.SIGTERM},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := unusedServerAddress(t)
			process := startServerProcess(t, address)
			waitForServer(t, process, address)

			require.NoError(t, process.cmd.Process.Signal(tt.signal))
			process.wait(t)
			require.NoError(t, process.waitErr, process.output.String())
			require.NotContains(t, process.output.String(), "panic:")
			require.NotContains(t, process.output.String(), "\ngoroutine ")

			connection, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
			if connection != nil {
				_ = connection.Close()
			}
			require.Error(t, err, "server must stop accepting connections")
		})
	}
}

type serverProcess struct {
	cmd     *exec.Cmd
	output  bytes.Buffer
	done    chan struct{}
	waitErr error
}

func startServerProcess(t *testing.T, address string) *serverProcess {
	t.Helper()

	process := &serverProcess{done: make(chan struct{})}
	process.cmd = exec.Command(os.Args[0], "-test.run=^TestServerProcessHelper$")
	process.cmd.Env = append(os.Environ(),
		serverHelperEnv+"=1",
		serverAddressEnv+"="+address,
	)

	process.cmd.Stdout = &process.output
	process.cmd.Stderr = &process.output
	require.NoError(t, process.cmd.Start())

	go func() {
		process.waitErr = process.cmd.Wait()
		close(process.done)
	}()

	t.Cleanup(func() {
		select {
		case <-process.done:
			return
		default:
			_ = process.cmd.Process.Kill()
		}
		select {
		case <-process.done:
		case <-time.After(processTimeout):
			t.Error("server process did not exit after cleanup")
		}
	})
	return process
}

func (process *serverProcess) wait(t *testing.T) {
	t.Helper()
	select {
	case <-process.done:
	case <-time.After(processTimeout):
		t.Fatal("server process did not exit before the deadline")
	}
}

func unusedServerAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	return address
}

func waitForServer(t *testing.T, process *serverProcess, address string) {
	t.Helper()
	client := &http.Client{
		Timeout: 200 * time.Millisecond,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	defer client.CloseIdleConnections()
	deadline := time.Now().Add(processTimeout)

	for time.Now().Before(deadline) {
		select {
		case <-process.done:
			t.Fatalf("server exited before readiness: %v\n%s", process.waitErr, process.output.String())
		default:
		}

		response, err := client.Get("http://" + address + "/")
		if err == nil {
			body, readErr := io.ReadAll(response.Body)
			_ = response.Body.Close()
			if readErr == nil && response.StatusCode == http.StatusOK &&
				strings.Contains(string(body), "<h1>Метрики</h1>") {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("server did not become ready before the deadline")
}
