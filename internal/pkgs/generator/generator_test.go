// Package generator provides generators for different purposes.
package generator

import "testing"

func TestIsUUID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{s: ""},
			want: false,
		},
		{
			name: "rand",
			args: args{s: "sldjfalsj"},
			want: false,
		},
		{
			name: "Uppercase",
			args: args{s: "C67DAB37-7ED3-5AA9-BD68-43A8FE0755B4"},
			want: false,
		},
		{
			name: "urn",
			args: args{s: "urn:uuid:123e4567-e89b-12d3-a456-426655440000"},
			want: false,
		},
		{
			name: "microsoft",
			args: args{s: "{123e4567-e89b-12d3-a456-426655440000}"},
			want: false,
		},
		{
			name: "ok",
			args: args{s: "f29e8857-2ffb-5871-9007-ee4de08f16c3"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUUID(tt.args.s); got != tt.want {
				t.Errorf("IsUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
