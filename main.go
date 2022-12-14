package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	//slsa02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/rhysd/actionlint"
	"github.com/slsa-framework/slsa-github-generator/slsa"

	"github.com/AdamKorcz/java-slsa-generator/gh"
)

// flags
var (
	flagYmlFile = flag.String("ymlfile", "", "path to ymlfile to attest")
)

var (
	demoBuildDir = "/home/runner/work/AdamKorcz/actions-test"
	demoBuilder  = "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_java_slsa3.yml@v1.4.0"
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

// Input: actionlint.Workflow.Jobs
// Output: The deploy command
func getSteps(jobs map[string]*actionlint.Job) []SlsaStep {
	provenanceSteps := make([]SlsaStep, 0)
	for _, job := range jobs {
		if len(job.Steps) != 0 {
			for _, step := range job.Steps {
				if step.Name != nil {
					//fmt.Println("Step name: ", step.Name.Value)
					//fmt.Println("Step exec: ", step.Exec.Kind())
					switch step.Exec.Kind() {
					case 0:
					case 1:
						if runCmd := step.Exec.(*actionlint.ExecRun).Run; runCmd != nil {
							runCmdSlice := strings.Split(runCmd.Value, " ")
							if len(runCmdSlice) == 0 {
								continue
							}
							if len(runCmdSlice) != 1 && runCmdSlice[len(runCmdSlice)-1] == "deploy" {
								break
							}
							if len(runCmdSlice) != 1 && runCmdSlice[len(runCmdSlice)-1] == "test" {
								break
							}
							provenanceStep := SlsaStep{
								Command:    strings.Split(runCmd.Value, " "),
								WorkingDir: demoBuildDir,
								Env:        []string{"env"},
							}
							provenanceSteps = append(provenanceSteps, provenanceStep)
						}
					}
				}
			}
		}
	}
	return provenanceSteps
}

func main() {
	flag.Parse()

	if *flagYmlFile == "" {
		fmt.Println("\nPlease specify a .yml file")
		fmt.Println("The usage is: go run main.go -ymlfile=/path/to/.yml\n")
		os.Exit(0)
	}

	b, err := os.ReadFile(*flagYmlFile)
	if err != nil {
		panic(err)
	}
	w, ymlErr := actionlint.Parse(b)
	if err != nil {
		panic(ymlErr)
	}

	testBuild := gh.NewTestBuild()
	g := slsa.NewHostedActionsGenerator(testBuild).WithClients(&slsa.NilClientProvider{})
	provenance, err := g.Generate(context.Background())
	if err != nil {
		panic(err)
	}
	steps := getSteps(w.Jobs)
	provenance.Predicate.BuildConfig = NewBuildConfig(steps)

	provenanceJson, err := json.MarshalIndent(provenance, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println("provenance: ", string(provenanceJson))
}
