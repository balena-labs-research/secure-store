package mtls

import (
	"testing"
)

func TestGenerateMTLSKeys(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 string
	}{
		{
			name:  "Check for Output",
			want:  "-----BEGIN CERTIFICATE-----",
			want1: "-----BEGIN RSA PRIVATE KEY-----",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GenerateMTLSKeys()
			if got[:len(tt.want)] != tt.want {
				t.Errorf("ValidateMTLSKeys() got = %v, want %v", got, tt.want)
			}
			if got1[:len(tt.want1)] != tt.want1 {
				t.Errorf("ValidateMTLSKeys() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateMTLSKeys(t *testing.T) {
	type args struct {
		encryptCert string
		encryptKey  string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Check for keys and cert",
			args: args{
				encryptCert: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZHRENDQXdDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0" +
					"ZBREFlTVJ3d0dnWURWUVFERXhOelpXTjEKY21VdGMzUnZjbVV0YzJWeWRtVnlNQ0FYRFRJeU1URXlNakUxTXpNeE0xb1lEekl4T" +
					"WpJeE1ESTVNVFV6TXpFegpXakFlTVJ3d0dnWURWUVFERXhOelpXTjFjbVV0YzNSdmNtVXRjMlZ5ZG1WeU1JSUNJakFOQmdrcWhr" +
					"aUc5dzBCCkFRRUZBQU9DQWc4QU1JSUNDZ0tDQWdFQWtzam5LdWhVMk8yUU5ZRnV6dDhMSjduczVuU3EvQ1JMelVZejV0TFkKZW1" +
					"sdGpRR2o2UDNOQndtcWs4czd5VWNDK25VWUJyQlVIVnVkTnE1a1hEWE80R2FmZkExS0pEVmxrQjE4cHNhLwpjbjdrczFUYkdLSE" +
					"ExNEVxbFBraE9Sbzl1dER0QVpIdG92K0xrQlBML3RLOHRWMUJpRjBPTFRZNDdaaWRSZGJ6CjJlbm9Ld2JDWVh5a2tQUlU2WkZmN" +
					"zBxbUtJRnBVMlhNaExpdWFVZ1BDMEpJZkJkdDhBWktVUWNlaVBtMVZ2b3UKWlFQS0ZkK3MvZ3N5emVpM2NtU0c0YmI5Vm5Ia09S" +
					"T1cwV0VZNWZjODM1MkoycE50Skp5ODZIVENQUkV6N3NCRgpyNytYRmR3UUpZd3c1V1hKWE5EU0dkVkZjT0hYcGFacy9nc2lLWFF" +
					"LODlMMmltWEZSRm85Q3VTa0F0Q05RUERjCkZOeXhYMHIvMVJkVlpzV0VNWHZuQjI4QXN4Mlp6UUpUTTQ4UjJQYlFjejJWN25GeH" +
					"ZGN3JRUFlmUjcxVHNaZ1EKZXVyU0VBZmxGOFZ4NWI5ZXBNRy9EeU1oSUorL08za2hUd1hsall2clluUFBZNFM3elEwMlRMUGVSW" +
					"XdCYVIzYQpGSlJ0dUNrMTZXTWI4ZXVDbElRbG5NNXhHVlFHSWVNLzBITnMwb29rdnhzZG41bGF2ZHpsc2hQUXQyRmkxZC9PCkQz" +
					"b3ZRbyt2aG9aMHNNblVWSzZpZ1BPUEFCSE1tbUdLMFV1VlhNQzFGanU4dzFmRDF3TlNCN1JJbXY0THYrT1gKRm1QRGk5N1FGYjV" +
					"6ZlNxREpjRUgrNzhLc25Zd2g0aTRjUWVTa2tveUMyT1FVdXAydmdQWHJQVkNacnM0Y0FBSApPMGNDQXdFQUFhTmZNRjB3RGdZRF" +
					"ZSMFBBUUgvQkFRREFnTzRNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CCkJnZ3JCZ0VGQlFjREFqQU1CZ05WSFJNQkFmOEVBa" +
					"kFBTUI0R0ExVWRFUVFYTUJXQ0UzTmxZM1Z5WlMxemRHOXkKWlMxelpYSjJaWEl3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUVl" +
					"V0J0OWV0bm1VVTI4N0N0OHp5NWo1OUhTbQpaSGdnL1lPck1ya0tHc1lWQkFUbjZhSnYrbXJaY2diYW5XNDgwT1FmSGtQSHRnR1p" +
					"MbjhzTjZsZkc2MVM5dWcyCkNFU3cvQ0pjTHUyN0tzRXZtaVFxVzdTeGlDTStqQ1oyZmZURHBwSFBUa2U1MVBiOXRjMm5ZdjNrUn" +
					"FRbkIxTncKeUFYNmNIdGpIVnpWYng0ZENjUUlrVUR0Y05Qc1JvTHQ2S1dOLzFmUjFseVZ0SWdwcXZTR2llVDR3dTdnVnZUTAp1S" +
					"VpKYW9NODVsYkpORUxROHhoMVR3UkpqZ041YVl1a1BuMnAyZllSUi8yN21EaW12bUlvZTlGRXY0RmVITDF3Ck14bzEySTJybU1L" +
					"R0cxUmd5M3hPbEVCUUFWLzRwbC85RTdQRVVTTG1kTWptV0lSdldKWHhQdllQRDZPSnd5L1QKbklXV0I0RzdhYy9LMXlMYVU5LzB" +
					"PQzgrZnVmelhHTTlRVm1GakI2OVAybTlKcjF2bm5jVmZYYjU5U3kvNWY4NwpQWFJYTWFtMXZ3RW5SRGpIYWpqcmpzbzNrSmgzSH" +
					"UxRmNYODJLTkVrSTlWMG0ycWI5R0FZczBiVnNKeDhOLzB2CmgzS25oNFdGR3V3TkwzUVM0NlR1UDVMcG16OXRxV1RxUzFnMzlHV" +
					"WF5TE1OVm0wMXEwMXI2a2Z3Z0haemxlTnIKcFd1T0tQQWtoUVh0ZW9iS3ZDdGhTazl3cG5CUVI2YXpHMWtjLzlyQWtvR1hyK0ZD" +
					"V1hiLzFCVEkzODhRNmRrZApDWVhZVlVDM1JUbTJpaGhHM2hzM2ZFbzRLQmVQSWtYZStsQ1FMRTNTSm1zSlBuWHphQmNWank1eE0" +
					"0SGI0WS9uCnVaaW1EeExyeHlQNjhhbEoKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
				encryptKey: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS1FJQkFBS0NBZ0VBa3Nqbkt1aFUyTzJRTllGdXp0OE" +
					"xKN25zNW5TcS9DUkx6VVl6NXRMWWVtbHRqUUdqCjZQM05Cd21xazhzN3lVY0MrblVZQnJCVUhWdWROcTVrWERYTzRHYWZmQTFLS" +
					"kRWbGtCMThwc2EvY243a3MxVGIKR0tIQTE0RXFsUGtoT1JvOXV0RHRBWkh0b3YrTGtCUEwvdEs4dFYxQmlGME9MVFk0N1ppZFJk" +
					"YnoyZW5vS3diQwpZWHlra1BSVTZaRmY3MHFtS0lGcFUyWE1oTGl1YVVnUEMwSklmQmR0OEFaS1VRY2VpUG0xVnZvdVpRUEtGZCt" +
					"zCi9nc3l6ZWkzY21TRzRiYjlWbkhrT1JPVzBXRVk1ZmM4MzUySjJwTnRKSnk4NkhUQ1BSRXo3c0JGcjcrWEZkd1EKSll3dzVXWE" +
					"pYTkRTR2RWRmNPSFhwYVpzL2dzaUtYUUs4OUwyaW1YRlJGbzlDdVNrQXRDTlFQRGNGTnl4WDByLwoxUmRWWnNXRU1Ydm5CMjhBc" +
					"3gyWnpRSlRNNDhSMlBiUWN6MlY3bkZ4dkY3clFQWWZSNzFUc1pnUWV1clNFQWZsCkY4Vng1YjllcE1HL0R5TWhJSisvTzNraFR3" +
					"WGxqWXZyWW5QUFk0Uzd6UTAyVExQZVJZd0JhUjNhRkpSdHVDazEKNldNYjhldUNsSVFsbk01eEdWUUdJZU0vMEhOczBvb2t2eHN" +
					"kbjVsYXZkemxzaFBRdDJGaTFkL09EM292UW8rdgpob1owc01uVVZLNmlnUE9QQUJITW1tR0swVXVWWE1DMUZqdTh3MWZEMXdOU0" +
					"I3UkltdjRMditPWEZtUERpOTdRCkZiNXpmU3FESmNFSCs3OEtzbll3aDRpNGNRZVNra295QzJPUVV1cDJ2Z1BYclBWQ1pyczRjQ" +
					"UFITzBjQ0F3RUEKQVFLQ0FnQVpVVjNPSU5UQnRmZ3h3bW1DZFNaUGE5cFl4YmJZVnNwY1ZjZU9BTUFtSFJrd1FTQXNoOXBkWkdW" +
					"dgpxSlpmV1VoQ093QXg1eWdiQ0RwTnZEYkRVT1NtQUExeU1EaWhsalEyYjErWXhKOHcrSUlxREhEZUJzaGtZM3NjCnh2dmNCWkt" +
					"3TG1reTVDREJCS0xsN3dRNHA5QzAxNmorQkU3MXJXeVVUc2FSelVteVJJNHZIYk5aYlU1VlVrOXoKR2dnYUZoNGl6bStFV3dLen" +
					"djY1RrY0pHaHVCSnh5REhUYzV6dTI2ZzRiRk5sMGw1VXZZQVg5eE8xVSttQklYZgp6VWNXbDZHS3o0bzE3aE1OZWxuNUVIUUhmZ" +
					"npSSHF5c3EveEpmQm9rdmxyZXkzczVuU3BLTXZMMkJ3QklFUTV5ClcvSmhVay92OElScHBoTVdub1Q5OFc0U01Fa2x3eGliVTBW" +
					"RVRZR2krNjlYUjExMjlVTFpFaUtaV0hnQ0twZVgKRjRUOCtLOXhQellSTVc2dDQ5RFpuR0tSQzdiRmNMK096SHFHczhTYk0zTER" +
					"DaXZPek9XOXFBR2ZGSUhGQnJ1TgpWWWpMdGlrK2Q1SkR3czJxZHRYV2s1M0RpVlBLMHlldDBwRE5VaksrMWpzZ3hKNjF3c3BwcF" +
					"p3bEFzbmpwcHN1CjNlaUd0UTlud2oxeUNhL29lY2liNEU1bnZaYnJLUXFVNHpqRnU3RllMSjdmT2ptOFpUZGQrRGJyRkR2QVRRQ" +
					"W4KdVJuVXBsUXpyV0gxdUFlUUVtRi9JTWo2U25WTjlSWitKVXpMYUMvRHo1ZU9ZazlwVEc1VENKSW1XY0FIdlFVRQpsOWJRSDNS" +
					"eStEdVRDZTBaWDJmWlQ1c3o3YmlBdHlZdWxIUnpMcFR3WU83UkRkRTJ3UUtDQVFFQXd2QXV1UUFvCjJRdDRPd2tNbjNMZ3h2SHp" +
					"aQ2drYTd4OWdwZ0F2ZHVpSnZxSHFZQVAxL2NDNGx1VlBMVmVjTEZja2J5aVh6WngKUUZlVFNWSEcwRmE5T0x5dGhhbmx5aDFYeG" +
					"NrWklyMTR2VkxaRHFXL1JKNDBlWnRsSFNybWpZUTlkZzZzcnlyWQphaUZOQzQ5RkNJekxKOVRnaWZ2U3BmMFZjUEZPYmo1VCtkR" +
					"mp4b0VjdnZFN0h1MDBZL2VOK24zc1J5RllSa0x0CmM1TGJNR0NzMjZ6L2N0eVNUaXJ5NW5HL1V2ZjdrRzUzR2d3eitpWjFJcENZ" +
					"U1RheEtkR1lyNHluQXlZV1Q3b3YKTUFNZEJBaTZyd09GcjdacXhMNFA5bzJVdVNtazlhY0dpU2RvcmJnbWt6RmZBUStTZG9Bd3F" +
					"pSTZCTFNvR2dNdAp4WFZTTmZCd0FuMUNYUUtDQVFFQXdNTmRkK3FqcUlvM3FxbEhNSUxKZmhCRFl3dWJ4QVoyVm5NZ28xSEZXMD" +
					"lTCkV3R2lMb2V2KytuMUtPUGhrUnY3c05uajZ0ZFgvZjV4YXRCUE9IRS93MStiejV2Mm1VdWtiZUtpUGswYjhuZGsKV3VsbC9LM" +
					"lBmQ2lFbmg4VGtFQXZJOUVvZUNwQnJCMVZEMWtwVjg0TVZaVjdDb2hNQno3MjYwamNiNWxHV1VsYwpZS2JOaW5JTGF2Q0s3SUxp" +
					"ODF4SlZiZTQrb1FndmlnSUZQS3NWODlpVGtJaEl5YW5UVDI0NjdWSjN2QW1HL3VZCmEvTGtzK1dSQkNWaWhqZlpKTzF0emhWdm9" +
					"XaTViRWpZQ1BrWi80ckpnVlMwbW1mc2dsMFdXTWlHMkNvK2NjRFAKWWxJVHAzQXE3UEE3UTcxZW1VTURXSGQrTmRyWW5sQ2krT0" +
					"9WaVN4aDh3S0NBUUVBcFFxd3BCVE1nS1pEVHNnawpsY0t2S0FDazdvRkdCS2o3SWx3TEZMTWxJQmN3VUlPSjVuRE5VbFB6a0FpR" +
					"1FxY0hGcVE4WEp3OVdocExLdUkvCmp4aEE4QTVlWXpJcXlPbjY3QXVNYW1zOEZCNVdneDQxUjZVaURHdFNPbFdlQ09hVEdxYmw4" +
					"UkEzVmZPSHhXZXgKTE1IM0ptd0hCd3ZibG9rbFpCQUpOVEV6NitncHIyQ1VzOXlORDJ3STFUSThWSDNVVTR5WnJqTHYwcy9kSDZ" +
					"KWgo5MHNLbzNhM3I5S1JBQ1lCcy84Y1QzWVhCRWljb1FkNldKOVBMMzFNbFhsTTZpUU5LbzlPaUlXN2VjekkvaWRkCjF0c0dqbl" +
					"BGQlhsZERvTmtEVGlDM3pCZ1NqTEJYNEx2dGVNdzZqY0M2NjB3aG5xZjRHZ0xncXA2VUJrNlYweVAKcUt2Z1NRS0NBUUJuRVk4T" +
					"3FYVnNqc3NKQTQ0L0VBOXd5bjFkbmoxaGFSc3BTbEd2UWUvR3pQalRsTklGRXlRMApESHE3ckcxVnk2S2VnMExzdE13bDBVRVBu" +
					"ZVplQlovVWRYQzRaMnBKdVRwb1R2LzFWQUdYZWVNOUFRTTVRZ3d3CmxQR2tJbld4Y2NCQzVIbEJGMGNhSW16eTBmZmJMUjlITzd" +
					"BWWg1d3ZXNmxuVVFRYzM5WFBzL3dpNmMySk0wV0wKbnFhSm84cDUyV1NVNHJHVHNjWllKc3Q3ZCs1UzNWNXcvY2IrSnBMYUtDQ1" +
					"NWOUx4MDNKdEQvQUYrRS8xS05DTAphUG52VGpsYkYyRFpDbmt2MHZFaVJ5ak1VM1ZsUCs1U3F0b09uWGJHbGNOM0lYR0liNzRIM" +
					"E5LU2ZkUm9lUDlhCnhYREszWnlSVWxXZG5mYnNFT0kyZEdOUHRWQ2xrZ2VsQW9JQkFRQzQ2MUNqUVlySFNSMSszT1RyZTNzYlJz" +
					"TncKUzVvSjM5L0RLVko3Zk9LYTc2VzJld2hVVW43bU5xelI1N2c3aTlibW5sZjdQbnJhOUkrN2txcFErWXlBaldzego0bDlFV3R" +
					"0MXBud0crQ21uZFVscTByNTdqdjgxMGdqNGNwTUNaM0pGZE5jWElmZ3VTY3RrMXpOaVBHL0RaRXVKCkVyS3pYMEtqNm1QQXRVWX" +
					"BQOVNiaXZmcS9kN0ZiU0plTmJSQnJ0Z3A5VlMzRE1uZ1FFZ0FGZjRlb1ovRGtSMnAKRkdxR3l3aUdXSnh1NmZPeVdaeWFaZWJtd" +
					"jJCWWcrT25IeEpnMzhIQit5SWNpTnlSZVhDZXBTZThWUXVkL2pJTApBa2dFb3JKdGF0Q2xlenZma0E3ekRIcndrU2xsa016WmI3" +
					"and2ai9nWUlMcHh6RGVuWXdqOVoxYU5EeFMKLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K",
			},
			want:  "-----BEGIN CERTIFICATE-----",
			want1: "-----BEGIN RSA PRIVATE KEY-----",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ValidateMTLSKeys(tt.args.encryptCert, tt.args.encryptKey)
			if got[:len(tt.want)] != tt.want {
				t.Errorf("ValidateMTLSKeys() got = %v, want %v", got, tt.want)
			}
			if got1[:len(tt.want1)] != tt.want1 {
				t.Errorf("ValidateMTLSKeys() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
