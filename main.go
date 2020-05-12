package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type cmDefinition struct {
	ConfigMaps []struct {
		Name  string
		Files []string
	} `yaml:"configMaps"`
}

type kubernetesCM struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Data       map[string]string
	Metadata   struct {
		Name string
	}
}

func ParseConfig(fileName string) (*cmDefinition, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	conf := &cmDefinition{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func BuildCM(conf *cmDefinition) (*[]*kubernetesCM, error) {
	var manifests []*kubernetesCM
	for _, cm := range conf.ConfigMaps {
		m := getEmptyManifest()
		m.Metadata.Name = cm.Name
		for _, fileName := range cm.Files {
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return nil, err
			}
			m.Data[filepath.Base(fileName)] = string(data)
		}
		manifests = append(manifests, m)
	}
	return &manifests, nil
}

func PrintOutput(manifests *[]*kubernetesCM) error {
	var outputs []string
	for _, m := range *manifests {
		out, err := yaml.Marshal(m)
		if err != nil {
			return err
		}
		outputs = append(outputs, string(out))
	}
	for _, o := range outputs {
		fmt.Printf("---\n%+v", o)
	}
	return nil
}

func getEmptyManifest() *kubernetesCM {
	return &kubernetesCM{
		ApiVersion: "v1",
		Kind:       "ConfigMap",
		Data:       make(map[string]string),
	}
}

func main() {
	const fn = ".file2cm.yaml"
	conf, err := ParseConfig(fn)
	if err != nil {
		log.Fatalf("Failed to read config from \"%v\": %v", fn, err)
	}
	mnf, err := BuildCM(conf)
	if err != nil {
		log.Fatalf("Failed to build manifest: %v", err)
	}
	err = PrintOutput(mnf)
	if err != nil {
		log.Fatalf("Failed to print manifest: %v", err)
	}
}
