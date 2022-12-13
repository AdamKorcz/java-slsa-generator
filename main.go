package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/rhysd/actionlint"
)

// flags
var (
	flagYmlFile = flag.String("ymlfile", "", "path to ymlfile to attest")
)


var (
	demoBuildDir = "/home/runner/work/AdamKorcz/actions-test"
	demoBuilder = "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/builder_java_slsa3.yml@v1.4.0"
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
	steps []*SlsaStep
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
	newStep := &SlsaStep{
		command: strings.Split(cmd, " "),

	}
	b.steps = append(b.steps, newStep)
}

// Creates a new buildConfig and adds the steps from the parsed workflow
func NewBuildConfig() *BuildConfig {
	b := &BuildConfig {
		steps: make([]*SlsaStep, 0),
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

type Predicate struct {
	Builder map[string]string
	BuildType string
	BuildConfig *BuildConfig
}

func NewPredicate() *Predicate {
	return &Predicate {
		Builder: map[string]string {
			"id": demoBuilder,
		},
		BuildType: "https://github.com/slsa-framework/slsa-github-generator/java@v1",
	}
}

func (p *Predicate) SetBuildConfig(b *BuildConfig) {
	p.BuildConfig = b
}

func main() {
	flag.Parse()

	if *flagYmlFile != "" {
		panic("A yml file must be specified")
	}

	b, err := os.ReadFile(flagYmlFile)
	if err != nil {
		panic(err)
	}
	w, ymlErr := actionlint.Parse(b)
	if err != nil {
		panic(ymlErr)
	}
	fmt.Printf("%+v\n", w)
	fmt.Println("Name: ", *w.Name)

	buildConf := NewBuildConfig()
	buildConf.SetSteps(w)

	p := NewPredicate()
	p.SetBuildConfig(buildConf)
	fmt.Println(p)
	
	fmt.Println("buildCon.steps = ", buildConf.steps[0].command)
}