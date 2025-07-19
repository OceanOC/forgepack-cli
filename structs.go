package main

type CFManifest struct {
	Minecraft       CFMinecraft `json:"minecraft"`
	ManifestType    string      `json:"manifestType"`
	ManifestVersion int         `json:"manifestVersion"`
	Name            string      `json:"name"`
	Version         string      `json:"version"`
	Author          string      `json:"author"`
	Files           []CFFile    `json:"files"`
}

type CFFile struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}

type CFMinecraft struct {
	Version        string       `json:"version"`
	Modloaders     []ModLoaders `json:"modLoaders"`
	ReccomendedRAM int          `json:"reccomendedRam"`
}

type ModLoaders struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}
