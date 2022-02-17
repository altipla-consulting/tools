package kubernetes

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScript(t *testing.T) {
	content, err := ioutil.ReadFile("test.jsonnet")
	require.NoError(t, err)
	buf, err := runScript("test", content)
	require.NoError(t, err)

	golden, err := ioutil.ReadFile("golden.json")
	require.NoError(t, err)
	require.JSONEq(t, string(golden), buf.String())
}
