package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"libs.altipla.consulting/errors"
)

func main() {
	if err := run(); err != nil {
		log.Error(err.Error())
		log.Debug(errors.Stack(err))
		os.Exit(1)
	}
}

func run() error {
	var flagDebug bool
	var flagStage, flagSource string
	flag.BoolVarP(&flagDebug, "debug", "d", false, "Enable debug logging")
	flag.StringVarP(&flagStage, "stage", "s", "", "Stage folder where the code will be transferred. If empty a temporary folder will be created.")
	flag.StringVarP(&flagSource, "source", "f", ".", "Source folder")
	flag.Parse()

	log.SetFormatter(new(log.TextFormatter))
	if flagDebug {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG log level activated")
	}

	if flagStage == "" {
		var err error
		flagStage, err = ioutil.TempDir("", "gaestage")
		if err != nil {
			return errors.Trace(err)
		}
	}

	log.WithFields(log.Fields{
		"source": flagSource,
		"stage":  flagStage,
	}).Info("Staging files to prepare them for deploy")

	exclude, err := ignore.CompileIgnoreFile(filepath.Join(flagSource, ".gcloudignore"))
	if err != nil {
		return errors.Trace(err)
	}

	var copied int64
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Trace(err)
		}

		rel, err := filepath.Rel(flagSource, path)
		if err != nil {
			return errors.Trace(err)
		}

		if info.IsDir() {
			// Ignore root folder, it is always included by definition.
			if rel == "." {
				return nil
			}

			if r, how := exclude.MatchesPathHow(rel); r {
				log.WithFields(log.Fields{
					"path": rel,
					"line": how.LineNo,
				}).Debug("Excluded folder")
				return filepath.SkipDir
			}

			return nil
		}

		if r, how := exclude.MatchesPathHow(rel); r {
			log.WithFields(log.Fields{
				"path": rel,
				"line": how.LineNo,
			}).Debug("Excluded file")
			return nil
		}

		destPath := filepath.Join(flagStage, rel)
		if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
			return errors.Trace(err)
		}
		log.WithFields(log.Fields{
			"path": rel,
			"dest": destPath,
		}).Debug("Copying file to staging directory")

		src, err := os.Open(path)
		if err != nil {
			return errors.Trace(err)
		}
		defer src.Close()
		dest, err := os.Create(destPath)
		if err != nil {
			return errors.Trace(err)
		}
		defer dest.Close()
		if _, err := io.Copy(dest, src); err != nil {
			return errors.Trace(err)
		}

		copied++

		return nil
	}
	if err := filepath.Walk(flagSource, fn); err != nil {
		return errors.Trace(err)
	}

	log.WithField("files", copied).Info("All files staged succesfully!")

	if flag.NArg() == 0 {
		return nil
	}

	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = flagStage
	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
