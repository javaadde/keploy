//go:build windows

package utils

import (
	"testing"
)

func TestGetPermissionFixCommand(t *testing.T) {
	tests := []struct {
		name       string
		keployPath string
		want       string
	}{
		{
			name:       "standard keploy path",
			keployPath: "C:\\Users\\user\\keploy",
			want:       "",
		},
		{
			name:       "empty path",
			keployPath: "",
			want:       "",
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
