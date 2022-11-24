package http

import (
	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mtls"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	type args struct {
		storePassword    string
		rcloneConfigPass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Start Server and Request the Key",
			args: args{
				storePassword:    "user-password",
				rcloneConfigPass: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mtls.GenerateMTLSKeys()
		flags.ServerHostname = "localhost"
		flags.Port = "8443"
		t.Run(tt.name, func(t *testing.T) {
			go StartServer(tt.args.storePassword, tt.args.rcloneConfigPass)
			time.Sleep(3 * time.Second)
			if err := StartClient(); (err != nil) != tt.wantErr {
				t.Errorf("StartServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
