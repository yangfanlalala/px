package slice

import (
	"fmt"
	"testing"
)

func TestReverse(t *testing.T) {
	type args struct {
		slice interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testing-01",
			args: args{
				slice: []uint32{},
			},
		},
		{
			name: "testing-02",
			args: args{
				slice: []uint32{1},
			},
		},
		{
			name: "testing-03",
			args: args{
				slice: []uint32{1, 2, 3, 4, 5, 6, 7, 8},
			},
		},
		{
			name: "testing-04",
			args: args{
				slice: &[]uint32{1, 2, 3, 4, 5, 6, 7, 8},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reverse(tt.args.slice)
			fmt.Println(tt.args.slice)
		})
	}
}
