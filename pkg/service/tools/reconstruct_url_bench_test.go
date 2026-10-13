package tools

import (
	"testing"
)

func TestReconstructURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		segmentPath string
		template    string
		want        string
	}{
		{
			name:        "simple url replace path segment 1",
			url:         "http://localhost:8080/users/123/profile",
			segmentPath: "path.1",
			template:    "{{string .id}}",
			want:        "http://localhost:8080/users/{{string .id}}/profile",
		},
		{
			name:        "url with query params",
			url:         "http://example.com/api/v1/resource?query=test&page=1",
			segmentPath: "path.2",
			template:    "{{string .resId}}",
			want:        "http://example.com/api/v1/{{string .resId}}?query=test&page=1",
		},
		{
			name:        "url with port and multiple segments",
			url:         "https://api.example.com:8443/a/b/c?x=1",
			segmentPath: "path.0",
			template:    "{{string .first}}",
			want:        "https://api.example.com:8443/{{string .first}}/b/c?x=1",
		},
		{
			name:        "out of bounds segment index - no change",
			url:         "http://localhost:8080/users/123",
			segmentPath: "path.5",
			template:    "{{string .id}}",
			want:        "http://localhost:8080/users/123",
		},
		{
			name:        "invalid segment path - no change",
			url:         "http://localhost:8080/users/123",
			segmentPath: "invalid",
			template:    "{{string .id}}",
			want:        "http://localhost:8080/users/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.url
			reconstructURL(&got, tt.segmentPath, tt.template)
			if got != tt.want {
				t.Errorf("reconstructURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkReconstructURL(b *testing.B) {
	urlStr := "http://example.com:8080/api/v1/users/12345/details?query=test&page=1"
	segmentPath := "path.3"
	template := "{{string .userId}}"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		u := urlStr
		reconstructURL(&u, segmentPath, template)
	}
}
