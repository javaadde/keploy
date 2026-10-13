//go:build !windows

package utils

import (
	"syscall"
	"testing"
)

func TestSetAndRestoreUmask_Others(t *testing.T) {
	// Get current umask
	current := syscall.Umask(022)
	syscall.Umask(current) // restore initial

	// Test SetUmask sets umask to 0 and returns previous mask
	oldMask := SetUmask()
	if oldMask != current {
		t.Errorf("expected SetUmask() to return previous umask 0%o, got 0%o", current, oldMask)
	}

	// Verify current umask was set to 0
	checkZero := syscall.Umask(022)
	if checkZero != 0 {
		t.Errorf("expected umask to be 0 after SetUmask(), got 0%o", checkZero)
	}

	// Test RestoreUmask restores the original mask
	RestoreUmask(oldMask)

	restored := syscall.Umask(022)
	syscall.Umask(restored) // restore
	if restored != oldMask {
		t.Errorf("expected RestoreUmask() to restore umask to 0%o, got 0%o", oldMask, restored)
	}
}
