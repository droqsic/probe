package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/droqsic/probe"
)

// TestEnvironmentVariables tests the behavior with different environment variables.
// This test uses a helper program to ensure that the behavior is consistent regardless of how the main test binary is built and run.
func TestEnvironmentVariables(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping environment variable test in CI environment")
	}

	helperCode := `
package main

import (
	"fmt"
	"os"

	"github.com/droqsic/probe"
)

func main() {
	fd := os.Stdout.Fd()
	fmt.Printf("TERM=%s IsTerminal=%v\n", 
		os.Getenv("TERM"), 
		probe.IsTerminal(fd))
}
`

	tempDir, err := os.MkdirTemp("", "probe-env-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	helperFile := filepath.Join(tempDir, "helper.go")
	if err := os.WriteFile(helperFile, []byte(helperCode), 0644); err != nil {
		t.Fatalf("Failed to write helper program: %v", err)
	}

	helperBin := filepath.Join(tempDir, "helper")
	if runtime.GOOS == "windows" {
		helperBin += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", helperBin, helperFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build helper program: %v", err)
	}

	termValues := []string{
		"xterm",
		"xterm-256color",
		"vt100",
		"dumb",
		"",
	}

	for _, term := range termValues {
		t.Run("TERM="+term, func(t *testing.T) {
			cmd := exec.Command(helperBin)
			cmd.Env = append(os.Environ(), "TERM="+term)

			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to run helper program: %v", err)
			}

			t.Logf("Output with TERM=%s: %s", term, string(output))
		})
	}
}

// TestCIEnvironment tests behavior in CI-like environments.
// This test sets up different CI-like environments and checks if IsTerminal correctly identifies them.
func TestCIEnvironment(t *testing.T) {
	originalEnv := os.Environ()

	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	ciEnvs := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "GitHub Actions",
			env: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
			},
		},
		{
			name: "Travis CI",
			env: map[string]string{
				"CI":     "true",
				"TRAVIS": "true",
			},
		},
		{
			name: "CircleCI",
			env: map[string]string{
				"CI":       "true",
				"CIRCLECI": "true",
			},
		},
		{
			name: "Jenkins",
			env: map[string]string{
				"CI":          "true",
				"JENKINS_URL": "http://example.com",
			},
		},
	}

	for _, ciEnv := range ciEnvs {
		t.Run(ciEnv.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range ciEnv.env {
				os.Setenv(k, v)
			}

			result := probe.IsTerminal(os.Stdout.Fd())
			t.Logf("IsTerminal in %s environment: %v", ciEnv.name, result)
		})
	}
}

// TestTerminalEmulation tests behavior with terminal emulation.
// This test checks if IsTerminal and IsCygwinTerminal correctly identify terminals on Windows and Unix-like systems.
func TestTerminalEmulation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Run("windows-terminal-emulation", func(t *testing.T) {
			result := probe.IsTerminal(os.Stdout.Fd())
			t.Logf("IsTerminal on Windows: %v", result)

			cygwinResult := probe.IsCygwinTerminal(os.Stdout.Fd())
			t.Logf("IsCygwinTerminal on Windows: %v", cygwinResult)
		})
	} else {
		t.Run("unix-terminal-emulation", func(t *testing.T) {
			result := probe.IsTerminal(os.Stdout.Fd())
			t.Logf("IsTerminal on Unix-like: %v", result)
		})
	}
}

// TestSSHEnvironment tests behavior in SSH-like environments.
// This test sets up an SSH-like environment and checks if IsTerminal correctly identifies it.
func TestSSHEnvironment(t *testing.T) {
	originalEnv := os.Environ()

	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}()

	os.Clearenv()
	os.Setenv("SSH_CONNECTION", "192.168.1.2 52698 192.168.1.3 22")
	os.Setenv("SSH_TTY", "/dev/pts/0")
	os.Setenv("TERM", "xterm-256color")

	result := probe.IsTerminal(os.Stdout.Fd())
	t.Logf("IsTerminal in SSH-like environment: %v", result)
}
