package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"libs.altipla.consulting/errors"
)

var lib = []byte(`
{
  imageVersion:: std.native('imageVersion'),
}
`)

func main() {
	if err := run(); err != nil {
		log.Fatal(errors.Stack(err))
	}
}

func run() error {
	if err := os.MkdirAll("tmp", 0700); err != nil {
		return errors.Trace(err)
	}

	dir, err := ioutil.TempDir("", "jnet")
	if err != nil {
		return errors.Trace(err)
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "jnet.jsonnet"), lib, 0600); err != nil {
		return errors.Trace(err)
	}

	var flagOutput, flagFilter string
	var flagEnv, flagDirectories []string
	flag.StringVarP(&flagOutput, "output", "o", "-", "Output JSON file")
	flag.StringVarP(&flagFilter, "filter", "f", "", "Filter top level items")
	flag.StringSliceVarP(&flagEnv, "env", "e", nil, "Environment variables")
	flag.StringSliceVarP(&flagDirectories, "directory", "d", nil, "Directories to include")
	flag.Parse()

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: append(flagDirectories, ".", dir),
	})

	for _, v := range flagEnv {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			return errors.Errorf("malformed environment variable: %s", v)
		}

		vm.ExtVar(parts[0], parts[1])
	}

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "imageVersion",
		Params: []ast.Identifier{"name"},
		Func: func(args []interface{}) (interface{}, error) {
			if os.Getenv("PROJECT_ID") != "" {
				return os.Getenv("BUILD_ID"), nil
			}

			cmd := exec.Command("docker", "image", "inspect", args[0].(string), "-f", "{{.Id}}")
			output, err := cmd.CombinedOutput()
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
					return "latest", nil
				}

				return nil, errors.Wrapf(err, "output: %s", output)
			}

			return strings.Split(string(output), ":")[1][:12], nil
		},
	})

	content, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		return errors.Trace(err)
	}

	list := &k8sList{
		APIVersion: "v1",
		Kind:       "List",
	}

	if flagFilter != "" {
		output, err := vm.EvaluateSnippetMulti(flag.Arg(0), string(content))
		if err != nil {
			return errors.Trace(err)
		}

		if output[flagFilter] == "" {
			output[flagFilter] = "{}"
		}

		var result interface{}
		if err := json.Unmarshal([]byte(output[flagFilter]), &result); err != nil {
			return errors.Trace(err)
		}
		extractItems(list, result)
	} else {
		output, err := vm.EvaluateSnippet(flag.Arg(0), string(content))
		if err != nil {
			return errors.Trace(err)
		}

		var result interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return errors.Trace(err)
		}
		extractItems(list, result)
	}

	var writer io.Writer
	if flagOutput == "-" {
		writer = os.Stdout
	} else {
		dir := filepath.Dir(flagOutput)
		if dir != "." {
			if err := os.MkdirAll(dir, 0700); err != nil {
				return errors.Trace(err)
			}
		}

		f, err := os.Create(flagOutput)
		if err != nil {
			return errors.Trace(err)
		}
		defer f.Close()

		writer = f
	}
	if err := json.NewEncoder(writer).Encode(list); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type k8sList struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Items      []interface{} `json:"items"`
}

func extractItems(list *k8sList, v interface{}) {
	switch v := v.(type) {
	case []interface{}:
		for _, x := range v {
			extractItems(list, x)
		}

	case map[string]interface{}:
		if _, ok := v["apiVersion"]; ok {
			list.Items = append(list.Items, v)
		} else {
			for _, x := range v {
				extractItems(list, x)
			}
		}

	case nil:
		return

	default:
		panic(fmt.Sprintf("should not reach here: %#v", v))
	}
}
