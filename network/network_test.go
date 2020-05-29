package network

import (
	"fmt"
	"testing"
)

func TestGenerateIPsByCIDR(t *testing.T) {
	type args struct {
		cidr string
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test01",
			args: args{cidr: "172.10.34.244/16"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIps, err := CIDRToIPRange(tt.args.cidr)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(len(gotIps))
			for _, ip := range gotIps {
				fmt.Println(ip)
			}
		})
	}
}