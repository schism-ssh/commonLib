package protocol

import "testing"

func TestGenerateLookupKey(t *testing.T) {
	type args struct {
		ident      string
		principals []string
	}
	tests := []struct {
		name string
		args args
		want LookupKey
	}{
		{
			name: "returns a lookup key for a host",
			args: args{
				ident:      "test.example.com",
				principals: []string{"test.example.com"},
			},
			want: LookupKey("55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"),
		},
		{
			name: "returns a lookup key for a user",
			args: args{
				ident:      "someUser@dev1.example.com",
				principals: []string{"someUser", "admin"},
			},
			want: LookupKey("a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateLookupKey(tt.args.ident, tt.args.principals); got != tt.want {
				t.Errorf("GenerateLookupKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
