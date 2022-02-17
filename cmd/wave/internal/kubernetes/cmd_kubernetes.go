package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/atlassian/go-sentry-api"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/wave/embed"
)

var (
	flagFilter   string
	flagEnv      []string
	flagIncludes []string
	flagApply    bool
)

func init() {
	Cmd.PersistentFlags().StringVarP(&flagFilter, "filter", "f", "", "Filter top level items when generating items.")
	Cmd.PersistentFlags().StringSliceVarP(&flagEnv, "env", "e", nil, "Set external variables.")
	Cmd.PersistentFlags().StringSliceVarP(&flagIncludes, "include", "i", nil, "Directories to include when running the jsonnet script.")
	Cmd.PersistentFlags().BoolVar(&flagApply, "apply", false, "Apply the output to the Kubernetes cluster instead of printing it.")
}

var Cmd = &cobra.Command{
	Use:     "kubernetes",
	Short:   "Run a jsonnet script and deploy the result to Kubernetes.",
	Example: "wave kubernetes k8s/deploy.jsonnet",
	Args:    cobra.ExactArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		sentryClient, err := sentry.NewClient(os.Getenv("SENTRY_AUTH_TOKEN"), nil, nil)
		if err != nil {
			return errors.Trace(err)
		}

		dir, err := ioutil.TempDir("", "jnet")
		if err != nil {
			return errors.Trace(err)
		}
		defer os.RemoveAll(dir)
		if err := ioutil.WriteFile(filepath.Join(dir, "wave.jsonnet"), embed.Wave, 0600); err != nil {
			return errors.Trace(err)
		}

		vm := jsonnet.MakeVM()
		vm.Importer(&jsonnet.FileImporter{
			JPaths: append(flagIncludes, ".", dir),
		})
		vm.NativeFunction(nativeFuncSentry(sentryClient))

		version := time.Now().Format("20060102") + "." + os.Getenv("BUILD_NUMBER")
		if ref := os.Getenv("GITHUB_REF"); ref != "" {
			version = path.Base(ref)
		}
		vm.ExtVar("version", version)

		for _, v := range flagEnv {
			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				return errors.Errorf("malformed environment variable: %s", v)
			}
			vm.ExtVar(parts[0], parts[1])
		}

		log.WithFields(log.Fields{
			"file":    args[0],
			"version": version,
		}).Info("Deploy generated file")

		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return errors.Trace(err)
		}

		list := &k8sList{
			APIVersion: "v1",
			Kind:       "List",
		}
		if flagFilter != "" {
			output, err := vm.EvaluateSnippetMulti(args[0], string(content))
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
			output, err := vm.EvaluateSnippet(args[0], string(content))
			if err != nil {
				return errors.Trace(err)
			}

			var result interface{}
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				return errors.Trace(err)
			}
			extractItems(list, result)
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(list); err != nil {
			return errors.Trace(err)
		}

		if !flagApply {
			fmt.Println(buf.String())
			return nil
		}

		apply := exec.Command("kubectl", "apply", "-f", "-")
		apply.Stdout = os.Stdout
		apply.Stderr = os.Stderr
		apply.Stdin = &buf
		if err := apply.Run(); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}

func apiString(s string) *string {
	return &s
}

func nativeFuncSentry(client *sentry.Client) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "sentry",
		Params: []ast.Identifier{"name"},
		Func: func(args []interface{}) (interface{}, error) {
			org := sentry.Organization{
				Slug: apiString("altipla-consulting"),
			}
			keys, err := client.GetClientKeys(org, sentry.Project{Slug: apiString(args[0].(string))})
			if err != nil {
				return nil, errors.Trace(err)
			}

			return keys[0].DSN.Public, nil
		},
	}
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
