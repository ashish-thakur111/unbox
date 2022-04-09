package utils

import (
	"log"
	"os"
	"text/template"
)

type FileParams struct {
	BaseImage string
	Volumes   []string
}

func ReadTmplAndDump(fileLocation string, params *FileParams) error {
	dir, err := createDir(".unbox")
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(fileLocation)
	if err != nil {
		return err
	}
	dockerFile, _ := os.Create(dir + "/Dockerfile")
	defer dockerFile.Close()
	err = tmpl.Execute(dockerFile, params)
	if err != nil {
		return err
	}
	return nil
}

/**
func readTmpl(fileLocation string) (string, error) {
	b, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
**/

func createDir(dirName string) (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := h + "/" + dirName
	if _, err := os.Stat(dir); err == nil {
		log.Println("work directory already exists, skipping creation...")
		return dir, nil
	}
	err = os.Mkdir(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}
