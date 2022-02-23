package kubernetes

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScript(t *testing.T) {
	os.Setenv("GITHUB_REF", "foo-version")

	content, err := ioutil.ReadFile("test.jsonnet")
	require.NoError(t, err)

	buf, err := runScript("test", content)
	require.NoError(t, err)

	fmt.Println(buf.String())

	golden, err := ioutil.ReadFile("golden.json")
	require.NoError(t, err)
	require.JSONEq(t, string(golden), buf.String())
}
