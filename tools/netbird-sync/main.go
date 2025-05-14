package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	config := getAppConfig()
	log.Printf("Config: %+v\n", config)
	log.Printf("Reading: %s\n", config.JSONPath)

	raw, err := os.ReadFile(config.JSONPath)
	if err != nil {
		panic(err)
	}

	var aclFile ACLFile
	err = json.Unmarshal(raw, &aclFile)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}

	// Parse resources
	policyResources := generatePolicyResources(aclFile.ACLs)
	groupResources := generateGroupResources(aclFile.Groups)

	// combine resources
	finalResources := map[string]any{}
	for k, v := range groupResources {
		finalResources[k] = v
	}
	for k, v := range policyResources {
		finalResources[k] = v
	}

	final := map[string]any{
		"resources": finalResources,
	}
	writeYAML(config, final)
}

func writeYAML(config *AppConfig, final any) error {
	// Use buffer with custom spaces
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(config.Indent)

	err := encoder.Encode(final)
	if err != nil {
		return fmt.Errorf("error encoding YAML: %v", err)
	}
	_ = encoder.Close()

	if config.OutPath != "" {
		// Save to output file
		err = os.WriteFile(config.OutPath, buf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("error writing to file %s: %v", config.OutPath, err)
		}
		log.Printf("YAML written to %s", config.OutPath)
	} else {
		fmt.Println(buf.String())
	}

	return nil
}
