package gh

import (
	"context"
	"github.com/slsa-framework/slsa-github-generator/github"
	"github.com/slsa-framework/slsa-github-generator/slsa"
)

var (
	TestBuildType   = "https://github.com/slsa-framework/slsa-github-generator/maven@v1"
	TestBuildConfig = "test build config"
)

type TestBuild struct {
	*slsa.GithubActionsBuild
}

func (*TestBuild) URI() string {
	return TestBuildType
}

func (*TestBuild) BuildConfig(context.Context) (interface{}, error) {
	return TestBuildConfig, nil
}

func NewTestBuild() *TestBuild {
	return &TestBuild{
		GithubActionsBuild: slsa.NewGithubActionsBuild(nil, github.WorkflowContext{
			RunID:      "12345",
			RunAttempt: "1",
			EventName:  "pull_request",
			SHA:        "abcde",
			RefType:    "branch",
			Ref:        "some/ref",
			BaseRef:    "some/base_ref",
			HeadRef:    "some/head_ref",
			RunNumber:  "102937",
			Actor:      "user",
		}).WithClients(&slsa.NilClientProvider{}),
	}
}
