/*
Copyright © 2022 Ashish Thakur

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"archive/zip"
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashish-thakur111/unbox/pkg/models"
	"github.com/ashish-thakur111/unbox/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// jarCmd represents the jar command
var jarCmd = &cobra.Command{
	Use:   "jar",
	Short: "unpacking and packing a jar file",
	Long:  `fat jar file to unpack and create a layered docker image file`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlLoc, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}
		if yamlLoc == "" {
			log.Panic("Please provide a yaml file")
		}
		config := DoReadYaml(yamlLoc)
		manifest, err := DoUnzipAndReadManifestfile(config)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range manifest {
			log.Println(k, v)
		}
		fileParams := utils.FileParams{config.Base, config.Context.Volumes}
		err = utils.ReadTmplAndDump("template/Dockerfile.tmpl", &fileParams)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func DoReadYaml(loc string) *models.Config {
	log.Println("Reading yaml filen from location" + loc)
	if !filepath.IsAbs(loc) {
		loc, _ = filepath.Abs(loc)
	}
	log.Println("reading yaml file")
	yamlFile, err := ioutil.ReadFile(loc)
	if err != nil {
		log.Fatal(err)
	}
	var config models.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}

func DoUnzipAndReadManifestfile(c *models.Config) (models.Manifest, error) {
	var jarPath string
	if filepath.IsAbs(c.Repo) {
		jarPath = c.Repo
	}
	u, err := url.ParseRequestURI(c.Repo)
	if err != nil {
		return nil, err
	}
	if u.Hostname() != "" {
		resp, err := http.Get(u.String())
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		defer resp.Body.Close()
		out, err := os.CreateTemp("", "fat-jar")
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}
		jarPath = out.Name()
	}
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, f := range r.File {
		if f.Name != "META-INF/MANIFEST.MF" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return readManifestData(rc)
	}
	defer r.Close()
	return nil, ErrNotJAR
}

var ErrNotJAR = errors.New("given file is not a JAR file")
var ErrWrongManifestFormat = errors.New("can't parse manifest file (wrong format)")

// readManifestData reads manifest data
func readManifestData(r io.Reader) (models.Manifest, error) {
	m := make(models.Manifest)
	s := bufio.NewScanner(r)

	var propName, propVal string

	for s.Scan() {
		text := s.Text()

		if len(text) == 0 {
			continue
		}

		if strings.HasPrefix(text, " ") {
			m[propName] += strings.TrimLeft(text, " ")
			continue
		}

		propSepIndex := strings.Index(text, ": ")

		if propSepIndex == -1 || len(text) < propSepIndex+2 {
			return nil, ErrWrongManifestFormat
		}

		propName = text[:propSepIndex]
		propVal = text[propSepIndex+2:]

		m[propName] = propVal
	}

	return m, nil
}

func init() {
	rootCmd.AddCommand(jarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	jarCmd.Flags().StringP("file", "f", "file", "yaml file to read")
}
