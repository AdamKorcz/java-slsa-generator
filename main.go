package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	slsa02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
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

func runToStep(runValue string) []string {
	return strings.Split(runValue, " ")
}

func isDeployStep(runValue string) bool {
	slice := strings.Split(runValue, " ")
	for _, elem := range slice {
		if elem == "deploy" {
			return true
		}
	}
	return false
}

type BuildConfig struct {
	version string
	steps   []SlsaStep
}

// Input: actionlint.Workflow.Jobs
// Output: The deploy command
func getDeployCmd(jobs map[string]*actionlint.Job) string {
	for _, job := range jobs {
		if len(job.Steps) != 0 {
			for _, step := range job.Steps {
				if step.Name != nil {
					//fmt.Println("Step name: ", step.Name.Value)
					//fmt.Println("Step exec: ", step.Exec.Kind())
					switch step.Exec.Kind() {
					case 0:
					case 1:
						if step.Exec.(*actionlint.ExecRun).Run != nil {
							if isDeployStep(step.Exec.(*actionlint.ExecRun).Run.Value) {
								return step.Exec.(*actionlint.ExecRun).Run.Value
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (b *BuildConfig) addStep(cmd string) {
	newStep := SlsaStep{
		command: strings.Split(cmd, " "),
	}
	b.steps = append(b.steps, newStep)
}

func totoSteps(w *actionlint.Workflow) [][]string {
	cmd := getDeployCmd(w.Jobs)

	steps := make([][]string, 0)

	if cmd != "" {
		steps = append(steps, strings.Split(cmd, " "))
	}
	return steps
}

// Creates a new buildConfig and adds the steps from the parsed workflow
func NewBuildConfig() *BuildConfig {
	b := &BuildConfig{
		steps:   make([]SlsaStep, 0),
		version: "1",
	}
	return b
}

func (b *BuildConfig) SetSteps(w *actionlint.Workflow) {
	dc := getDeployCmd(w.Jobs)

	if dc != "" {
		fmt.Println("deploy step: ", dc)
		b.addStep(dc)
	}

}

type SlsaStep struct {
	command    []string
	env        []string
	workingDir string
}

func SetTotoProvBuildConfig(predicate *slsa02.ProvenancePredicate, buildConfig interface{}) {
	predicate.BuildConfig = buildConfig
}

func main() {
	flag.Parse()

	if *flagYmlFile == "" {
		panic("A yml file must be specified")
	}

	b, err := os.ReadFile(*flagYmlFile)
	if err != nil {
		panic(err)
	}
	w, ymlErr := actionlint.Parse(b)
	if err != nil {
		panic(ymlErr)
	}
	_ = w

	testBuild := gh.NewTestBuild()
	g := slsa.NewHostedActionsGenerator(testBuild).WithClients(&slsa.NilClientProvider{})
	provenance, err := g.Generate(context.Background())
	if err != nil {
		panic(err)
	}
	steps := totoSteps(w)
	provenance.Predicate.BuildConfig = steps

	provenanceJson, err := json.MarshalIndent(provenance, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println("provenance: ", string(provenanceJson))
}
