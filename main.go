package main

import (
	"context"
	"encoding/json"
	"fmt"

	//slsa02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/slsa-framework/slsa-github-generator/slsa"

	"github.com/AdamKorcz/java-slsa-generator/gh"
)

type SlsaStep struct {
	Command    []string `json:"command"`
	Env        []string `json:"env"`
	WorkingDir string   `json:"workingDir"`
}

type BuildConfig struct {
	Version int        `json:"version"`
	Steps   []SlsaStep `json:"steps"`
}

func NewBuildConfig(steps []SlsaStep) BuildConfig {
	return BuildConfig{
		Version: 1,
		Steps:   steps,
	}
}

func main() {
	testBuild := gh.NewTestBuild()
	g := slsa.NewHostedActionsGenerator(testBuild).WithClients(&slsa.NilClientProvider{})
	provenance, err := g.Generate(context.Background())
	if err != nil {
		panic(err)
	}

	//provenance.Predicate.BuildConfig = NewBuildConfig(steps)

	provenanceJson, err := json.MarshalIndent(provenance, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println("provenance: ", string(provenanceJson))
}
