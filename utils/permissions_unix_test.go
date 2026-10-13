//go:build linux || darwin

package utils

import (
	"fmt"
	"os/user"
	"testing"
)

func TestGetPermissionFixCommand(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("failed to get current user: %v", err)
	}

	tests := []struct {
		name       string
		keployPath string
		want       string
	}{
		{
			name:       "standard keploy path",
			keployPath: "/home/user/keploy",
			want:       fmt.Sprintf("sudo chown -R %s /home/user/keploy", currentUser.Username),
		},
		{
			name:       "path with spaces",
			keployPath: "/home/user/my keploy path",
			want:       fmt.Sprintf("sudo chown -R %s /home/user/my keploy path", currentUser.Username),
		},
		{
			name:       "relative path",
			keployPath: "./keploy",
			want:       fmt.Sprintf("sudo chown -R %s ./keploy", currentUser.Username),
		},
		{
			name:       "empty path",
			keployPath: "",
			want:       fmt.Sprintf("sudo chown -R %s ", currentUser.Username),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(subT *testing.T) {
			got := GetPermissionFixCommand(tt.keployPath)
			if got != tt.want {
				subT.Errorf("GetPermissionFixCommand(%q) = %q, want %q", tt.keployPath, got, tt.want)
			}
		})
	}
}
