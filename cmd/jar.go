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
	"io/ioutil"
	"log"
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
