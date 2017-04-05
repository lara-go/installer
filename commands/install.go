package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

const boilerplateZipURL = "https://github.com/lara-go/boilerplate/archive/master.zip"

func installHelp(gopath string) {
	fmt.Printf(`
  Generate a new project from the boilerplate in %s

  Usage: larago install [project-name]
`, gopath)
	os.Exit(0)
}

// Install command.
func Install(args []string, verbose bool) {
	gopath := os.Getenv("GOPATH")

	if len(args) == 1 {
		installHelp(gopath)
	}

	tmpDir := os.TempDir()
	project := args[1]
	zipFile := path.Join(tmpDir, "larago-boilerplate.zip")
	boilerplatePath := path.Join(tmpDir, "boilerplate-master")
	projectPath := path.Join(gopath, "src", project)

	if _, err := os.Stat(projectPath); err == nil {
		fmt.Println("Project already exists in", projectPath)
		os.Exit(0)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Downloading boilerplate: "
	s.FinalMSG = "\n"
	if !verbose {
		s.Start()
	}

	if verbose {
		fmt.Println("Downloading", boilerplateZipURL, "to", zipFile)
	}
	if err := downloadFromURL(boilerplateZipURL, zipFile); err != nil {
		showError(err)
	}

	if verbose {
		fmt.Println("Unzipping", boilerplateZipURL, "to", tmpDir)
	}
	if err := unzip(zipFile, tmpDir); err != nil {
		showError(err)
	}

	if verbose {
		fmt.Println("Moving boilerplate", boilerplatePath, "to", projectPath)
	}
	if err := os.Rename(boilerplatePath, projectPath); err != nil {
		showError(err)
	}

	if verbose {
		fmt.Println("Update imports")
	}
	err := filepath.Walk(projectPath, func(path string, fi os.FileInfo, err error) error {
		return updateImports(path, fi, err, project)
	})

	if err != nil {
		showError(err)
	}

	if !verbose {
		s.Stop()
	}

	fmt.Printf(`New project was installed in %s

Next steps:

  1. Install Glide tool (https://glide.sh)
  2. Install dependencies using Glide

To check, run:

  $ go run cmd/app/main.go -r ./ env
`, projectPath)
}

func showError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func downloadFromURL(url, dest string) error {
	// TODO: check file existence first with io.IsExist
	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if _, err := io.Copy(output, response.Body); err != nil {
		return err
	}

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if errr := f.Close(); errr != nil {
					panic(errr)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateImports(path string, fi os.FileInfo, err error, project string) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", fi.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		newContents := strings.Replace(
			string(read),
			"github.com/lara-go/boilerplate",
			project,
			-1,
		)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

	}

	return nil
}
