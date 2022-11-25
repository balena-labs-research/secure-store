//go:build linux

package mount

import (
	"os"
	"testing"
	"time"
)

func TestCreateMount(t *testing.T) {
	type args struct {
		passwd       string
		mountPoint   string
		storagePoint string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create an encrypted mount",
			args: args{
				passwd:       "my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ",
				mountPoint:   "./decrypted",
				storagePoint: "./storage",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateMount(tt.args.passwd); (err != nil) != tt.wantErr {
				t.Errorf("CreateMount() error = %v, wantErr %v", err, tt.wantErr)
			}

			f, err := os.Create(tt.args.mountPoint + "/filename.ext")
			if err != nil {
				t.Errorf("CreateMount() error = %v)", err)
			}
			f.Close()

			time.Sleep(2 * time.Second)

			_, err = os.Stat(tt.args.mountPoint)
			if os.IsNotExist(err) {
				t.Errorf("CreateMount() error = %v)", err)
			}

			_, err = os.Stat(tt.args.storagePoint)
			if os.IsNotExist(err) {
				t.Errorf("CreateMount() error = %v)", err)
			}

			unMount(tt.args.mountPoint)
		})
	}
}
