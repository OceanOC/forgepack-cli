package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func OpenZIP(file, outputFolder string, manifest *CFManifest) {
	outFolder := filepath.Clean(outputFolder + "/modpacks/")
	tempFolder := filepath.Clean(outputFolder + "/modpacks/temp/")
	os.MkdirAll(outFolder, 0755)
	os.MkdirAll(tempFolder, 0755)

	r, err := zip.OpenReader(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()
	fmt.Printf("Opened %s\n", file)

	for _, f := range r.File {
		if f.FileInfo().Name() == "manifest.json" {
			fmt.Println("Opened manifest.json")
			fo, err := f.Open()
			if err != nil {
				log.Fatalln(err)
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(fo)
			OpenManifest(buf.Bytes(), outFolder, manifest)

		} else if strings.HasPrefix(f.Name, "overrides/") {
			fmt.Printf("Extracting override: %s\n", filepath.Dir(tempFolder+"/"+f.Name))
			relPath := tempFolder + "/"
			err := os.MkdirAll(filepath.Dir(relPath+f.Name), 0755)
			if err != nil {
				log.Fatalln(err)
			}
			fo, err := f.Open()
			if err != nil {
				log.Fatalln(err)
			}
			file, err := os.Create(relPath + f.Name)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()

			_, err = io.Copy(file, fo)
			if err != nil {
				log.Fatalln(err)
			}

		}
	}

	modPackFolder := filepath.Clean(outFolder + "/" + strings.ReplaceAll(manifest.Name, " ", ""))

	err = os.CopyFS(modPackFolder, os.DirFS(tempFolder+"/overrides/"))
	if err != nil {
		log.Fatalln(err)
	}

	err = os.RemoveAll(tempFolder)
	if err != nil {
		log.Fatalln(err)
	}
}

func OpenManifest(file []byte, outFolder string, manifest *CFManifest) {
	err := json.Unmarshal(file, manifest)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("manifest.json parsed.")
	fmt.Printf("Modpack name: %s\n", manifest.Name)
	fmt.Printf("Modpack author: %s\n", manifest.Author)
	fmt.Printf("Modpack mod count: %d\n", len(manifest.Files))

	modPackFolder := filepath.Clean(outFolder + "/" + strings.ReplaceAll(manifest.Name, " ", ""))
	modFolder := filepath.Clean(modPackFolder + "/mods/")

	_, err = os.Stat(modPackFolder)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(modPackFolder, 0755)
			os.MkdirAll(modPackFolder+"/mods/", 0755)
		} else {
			log.Fatalln(err)
		}
	}

	var prog float32
	for i := 0; i < len(manifest.Files); i++ {
		out, err := os.Create(modFolder + "/" + strconv.Itoa(manifest.Files[i].FileID) + ".jar")
		if err != nil {
			log.Fatalln(err)
		}
		defer out.Close()

		fmt.Printf("Downloading mod: %s (%.2f%%, %d/%d) \n", strconv.Itoa(manifest.Files[i].FileID), prog, i, len(manifest.Files))

		resp, err := http.Get("https://www.curseforge.com/api/v1/mods/ " + strconv.Itoa(manifest.Files[i].ProjectID) + "/files/" + strconv.Itoa(manifest.Files[i].FileID) + "/download")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		prog = float32(i) / float32(len(manifest.Files)) * 100
	}
}
