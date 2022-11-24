package decrypt

import (
	"reflect"
	"testing"
)

func TestDecryptEnvs(t *testing.T) {
	type args struct {
		envs     []string
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check for Output",
			args: args{
				envs:     []string{"ENCRYPTED_TESTVAR=b2c8955bad9bc7520e09f42e7a24b34860745b6fad240368dcc201d2d175938f55b5a404265c09f848c758aeb4248b7b32d63d1baad66a"},
				password: "my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ",
			},
			want: "TESTVAR=this-is-my-new-test-api-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecryptEnvs(tt.args.envs, tt.args.password); !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("DecryptEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}
