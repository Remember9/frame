package xaes

import (
	"testing"
)

func TestAesCbcEncrypt(t *testing.T) {
	type args struct {
		plainTextOrg string
		keyOrg       string
		ivAesOrg     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestAesCbcEncrypt", args{"123456", "rDOMmS79sMsX9GwwKrMHxW6hkcHrFF9I", "z5UVKpjSbiE6FDzQ"}, "123456"},
		{"TestAesCbcEncrypt", args{"123456", "29759ba08f4b4e4caac36bf63a3e1374", "749029A541C64DD6"}, "123456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypt, err := AesCbcEncrypt(tt.args.plainTextOrg, tt.args.keyOrg, tt.args.ivAesOrg)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(encrypt)
			decrypt, err := AesCbcDecrypt(encrypt, tt.args.keyOrg, tt.args.ivAesOrg)
			if err != nil {
				t.Fatalf("AesCbcDecrypt() decrypt = %v, want %v", decrypt, tt.want)
			}
		})
	}
}
