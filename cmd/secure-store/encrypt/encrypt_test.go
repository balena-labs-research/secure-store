package encrypt

import "testing"

func TestEncryptString(t *testing.T) {
	type args struct {
		key  string
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Check Length",
			args: args{
				key:  "user-password",
				text: "test-text",
			},
			want: 74,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncryptString(tt.args.key, tt.args.text); len(got) != tt.want {
				t.Errorf("EncryptString() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGenerateNewKey(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Check Length",
			want: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateNewKey(); len(got) != tt.want {
				t.Errorf("GenerateNewKey() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
