package hash

import (
	"testing"
)

func TestHasher_IsValidSum(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		sum   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "is valid",
			fields: fields{
				key: "key",
			},
			args: args{
				sum:   "73ea8f3887ebbb6afdbc8a8be682faea0f22f82d02bf6fd84c469d6556babb63",
				value: "value",
			},
			want: true,
		},
		{
			name: "is not valid",
			fields: fields{
				key: "key",
			},
			args: args{
				sum:   "73ea8f3887ebbb6afdbc8a8be682faea0f22f82d02bf6fd84c469d6556babb63",
				value: "wrong value",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				h := Hasher{
					key: tt.fields.key,
				}
				if got := h.IsValidSum(tt.args.sum, tt.args.value); got != tt.want {
					t.Errorf("IsValidSum() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestHasher_Sha256(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Makes sha256 signature",
			fields: fields{
				key: "key",
			},
			args: args{
				value: "value",
			},
			want: "73ea8f3887ebbb6afdbc8a8be682faea0f22f82d02bf6fd84c469d6556babb63",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				h := Hasher{
					key: tt.fields.key,
				}
				if got := h.Sha256(tt.args.value); got != tt.want {
					t.Errorf("Sha256() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
