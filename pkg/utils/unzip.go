package utils

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
)

var ErrNotJAR = errors.New("given file is not a JAR file")
var ErrWrongManifestFormat = errors.New("can't parse manifest file (wrong format)")

// function to unzip jar and with manifest meta information
func UnzipJar(c *models.Config) (models.Manifest, string, error) {
	var jarPath string
	if filepath.IsAbs(c.Repo) {
		jarPath = c.Repo
	}
	u, err := url.ParseRequestURI(c.Repo)
	if err != nil {
		return nil, "", err
	}
	if u.Hostname() != "" {
		resp, err := http.Get(u.String())
		if err != nil {
			log.Fatalln(err)
			return nil, "", err
		}
		defer resp.Body.Close()
		out, err := os.CreateTemp("", "fat-jar")
		if err != nil {
			log.Fatalln(err)
			return nil, "", err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalln(err)
			return nil, "", err
		}
		jarPath = out.Name()
	}
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		log.Fatal(err)
		return nil, "", err
	}
	return extractJarAndReadManifest(r)
}

// function to extract jar and read manifest, return manifest and extracted location
func extractJarAndReadManifest(r *zip.ReadCloser) (models.Manifest, string, error) {
	dest, err := ioutil.TempDir("", "extracted-jar")
	if err != nil {
		log.Fatalln(err)
	}
	var md models.Manifest
	for _, f := range r.File {
		fp := filepath.Join(dest, f.Name)
		log.Println("unzipping file", fp)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		outFile, _ := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		rc, _ := f.Open()
		io.Copy(outFile, rc)

		defer outFile.Close()
		defer rc.Close()
		if f.Name == "META-INF/MANIFEST.MF" {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
				return nil, "", err
			}
			md, err = readManifestData(rc)
			if err != nil {
				log.Fatal(err)
				return nil, "", err
			}
		}
	}
	defer r.Close()
	return md, dest, nil
}

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
