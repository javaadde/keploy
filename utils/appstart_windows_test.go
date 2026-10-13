//go:build windows

package utils

import (
	"errors"
	"os/exec"
	"testing"
)

func TestSetNativeAppStarter(t *testing.T) {
	defer SetNativeAppStarter(nil)

	// Initially nativeAppStarter should be nil
	if fn := nativeAppStarter.Load(); fn != nil {
		t.Fatalf("expected nativeAppStarter to be nil initially, got non-nil")
	}

	var called bool
	dummyStarter := func(cmd *exec.Cmd) error {
		called = true
		return nil
	}

	// Set starter
	SetNativeAppStarter(dummyStarter)
	fnPtr := nativeAppStarter.Load()
	if fnPtr == nil {
		t.Fatalf("expected nativeAppStarter to be non-nil after SetNativeAppStarter")
	}

	// Invoke stored function to ensure it matches
	(*fnPtr)(nil)
	if !called {
		t.Fatalf("expected dummyStarter to be called")
	}

	// Overwrite starter with another function
	var calledSecond bool
	secondStarter := func(cmd *exec.Cmd) error {
		calledSecond = true
		return nil
	}
	SetNativeAppStarter(secondStarter)
	fnPtr2 := nativeAppStarter.Load()
	if fnPtr2 == nil {
		t.Fatalf("expected nativeAppStarter to be non-nil")
	}
	(*fnPtr2)(nil)
	if !calledSecond {
		t.Fatalf("expected secondStarter to be called")
	}

	// Clear starter
	SetNativeAppStarter(nil)
	if fn := nativeAppStarter.Load(); fn != nil {
		t.Fatalf("expected nativeAppStarter to be nil after passing nil to SetNativeAppStarter")
	}
}

func TestStartApplication(t *testing.T) {
	defer SetNativeAppStarter(nil)

	errCustom := errors.New("custom native starter called")

	t.Run("Native with starter set", func(t *testing.T) {
		defer SetNativeAppStarter(nil)

		var receivedCmd *exec.Cmd
		SetNativeAppStarter(func(cmd *exec.Cmd) error {
			receivedCmd = cmd
			return errCustom
		})

		testCmd := exec.Command("cmd.exe", "/c", "echo test")
		err := startApplication(testCmd, Native)
		if !errors.Is(err, errCustom) {
			t.Errorf("expected startApplication to return errCustom, got: %v", err)
		}
		if receivedCmd != testCmd {
			t.Errorf("expected custom starter to receive cmd %v, got %v", testCmd, receivedCmd)
		}
	})

	t.Run("Native without starter", func(t *testing.T) {
		defer SetNativeAppStarter(nil)
		SetNativeAppStarter(nil)

		testCmd := exec.Command("nonexistent_executable_12345")
		err := startApplication(testCmd, Native)
		if errors.Is(err, errCustom) {
			t.Errorf("expected startApplication NOT to call custom starter, got: %v", err)
		}
		if err == nil {
			t.Errorf("expected startApplication to return error from cmd.Start() on nonexistent executable")
		}
	})

	t.Run("Non-native with starter set", func(t *testing.T) {
		defer SetNativeAppStarter(nil)

		SetNativeAppStarter(func(cmd *exec.Cmd) error {
			return errCustom
		})

		for _, kind := range []CmdType{DockerRun, DockerStart, DockerCompose, Empty, FromContainer} {
			testCmd := exec.Command("nonexistent_executable_12345")
			err := startApplication(testCmd, kind)
			if errors.Is(err, errCustom) {
				t.Errorf("expected startApplication NOT to invoke custom starter for kind %q, but got errCustom", kind)
			}
			if err == nil {
				t.Errorf("expected startApplication to return error from cmd.Start() for kind %q", kind)
			}
		}
	})
}
