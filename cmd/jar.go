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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ashish-thakur111/unbox/pkg/models"
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
		DoUnzipAndCreateDockerfile(config)
	},
}

func DoReadYaml(yamlLoc string) *models.Config {
	log.Println("Reading yaml filen from location" + yamlLoc)
	if !filepath.IsAbs(yamlLoc) {
		yamlLoc, _ = filepath.Abs(yamlLoc)
	}
	log.Println("reading yaml file")
	yamlFile, err := ioutil.ReadFile(yamlLoc)
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

func DoUnzipAndCreateDockerfile(config *models.Config) {
	var jarPath string
	if filepath.IsAbs(config.Repo) {
		jarPath = config.Repo
	}
	u, err := url.ParseRequestURI(config.Repo)
	if err != nil {
		return
	}
	if u.Hostname() != "" {
		resp, err := http.Get(u.String())
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		out, err := os.CreateTemp("", "fat-jar")
		if err != nil {
			log.Fatalln(err)
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		jarPath = out.Name()
	}
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
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
