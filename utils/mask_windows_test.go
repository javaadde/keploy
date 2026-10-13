//go:build windows

package utils

import "testing"

func TestSetUmask_Windows(t *testing.T) {
	oldMask := SetUmask()
	if oldMask != 0 {
		t.Errorf("expected SetUmask() to return 0 on Windows, got %d", oldMask)
	}
}

func TestRestoreUmask_Windows(t *testing.T) {
	// RestoreUmask is a no-op on Windows.
	// Verify it executes without error/panic for any integer input.
	RestoreUmask(0)
	RestoreUmask(022)
	RestoreUmask(077)
}
