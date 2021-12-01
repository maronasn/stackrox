package maincommand

import (
	"testing"

	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/printer"
	"github.com/stretchr/testify/assert"
)

func TestErrorWriter(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "\nError: rpc error: code = Unauthenticated desc =\n credentials not found\n",
			out: "\nError: rpc error: code = Unauthenticated desc =\n credentials not found\n",
		},
		{
			in:  "rpc error: code = Unauthenticated desc = credentials not found",
			out: "rpc error: code = Unauthenticated desc = credentials not found\n",
		},
		{
			in:  "rpc error: code = Unauthenticated desc = credentials not found\n",
			out: "rpc error: code = Unauthenticated desc = credentials not found\n",
		},
		{
			in:  "Error: rpc error: code = Unauthenticated desc = credentials not found",
			out: "ERROR:\trpc error: code = Unauthenticated desc = credentials not found\n",
		},
		{
			in:  "",
			out: "ERROR:\t\n",
		},
		{
			in:  "%s",
			out: "ERROR:\t%s\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			io, _, out, errorOut := environment.TestIO()
			ew := errorWriter{
				logger: environment.NewLogger(io, printer.DefaultColorPrinter()),
			}
			n, err := ew.Write([]byte(tt.in))
			assert.NoError(t, err)
			assert.Len(t, tt.in, n)
			assert.Empty(t, out.String())
			assert.Equal(t, tt.out, errorOut.String())
		})
	}
}
