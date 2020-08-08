package mapping

import (
	"fmt"
	tu "github.com/skwiwel/site-mapper/test/testutil"
	"testing"
)

func Test_isInSearchScope(t *testing.T) {
	type args struct {
		address       fmt.Stringer
		masterAddress fmt.Stringer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "child URL",
			args: args{
				tu.URLParseSkipError("http://foo.com/users", t),
				tu.URLParseSkipError("http://foo.com", t),
			},
			want: true,
		},
		{
			name: "other domain",
			args: args{
				tu.URLParseSkipError("http://example.com", t),
				tu.URLParseSkipError("http://foo.com", t),
			},
			want: false,
		},
		{
			name: "address higher in hierarchy than master",
			args: args{
				tu.URLParseSkipError("http://foo.com/master", t),
				tu.URLParseSkipError("http://foo.com/master/users", t),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInSearchScope(tt.args.address, tt.args.masterAddress); got != tt.want {
				t.Errorf("isInSearchScope() = %v, want %v\naddress:       %s\nmasterAddress: %s", got, tt.want, tt.args.address.String(), tt.args.masterAddress.String())
			}
		})
	}
}
