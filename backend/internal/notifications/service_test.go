package notifications

import "testing"

func TestIsAllowedEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		want     bool
	}{
		{"https://fcm.googleapis.com/fcm/send/abc", true},
		{"https://updates.push.services.mozilla.com/wpush/v2/abc", true},
		{"https://web.push.apple.com/abc", true},
		{"https://wns2-db5p.notify.windows.com/w/?token=abc", true},
		{"http://fcm.googleapis.com/fcm/send/abc", false},
		{"https://fcm.googleapis.com:8443/fcm/send/abc", false},
		{"https://user@fcm.googleapis.com/fcm/send/abc", false},
		{"https://fcm.googleapis.com.evil.com/abc", false},
		{"https://evilfcm.googleapis.com.example/abc", false},
		{"https://169.254.169.254/computeMetadata/v1/", false},
		{"https://localhost/abc", false},
		{"not a url", false},
	}

	for _, tt := range tests {
		if got := isAllowedEndpoint(tt.endpoint); got != tt.want {
			t.Errorf("isAllowedEndpoint(%q) = %v, want %v", tt.endpoint, got, tt.want)
		}
	}
}
