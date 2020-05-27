// Package tristate provides a bool with 3 states.
package tristate

import (
	"testing"
)

func TestNewFromStr(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		{
			name: "t",
			args: args{
				s: "t",
			},
			want: True,
		},
		{
			name: "t",
			args: args{
				s: "t",
			},
			want: True,
		},
		{
			name: "asd",
			args: args{
				s: "asd",
			},
			want: NotSet,
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: NotSet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromStr(tt.args.s); got != tt.want {
				t.Errorf("NewFromStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	pos := true
	neg := false

	type args struct {
		b *bool
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		{
			name: "nil",
			args: args{b: nil},
			want: NotSet,
		},
		{
			name: "t",
			args: args{b: &pos},
			want: True,
		},
		{
			name: "f",
			args: args{b: &neg},
			want: False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.b); got != tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
