package main

import (
	"context"
	"testing"
	"github.com/google/go-cmp/cmp"
	"github.com/AdamKorcz/java-slsa-generator/gh"
	"github.com/slsa-framework/slsa-github-generator/slsa"
	slsa02 "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	slsacommon "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

func TestProvenance1(t *testing.T) {
	b := NewTestBuild()
	expected := &intoto.ProvenanceStatement{
				StatementHeader: intoto.StatementHeader{
					Type:          intoto.StatementInTotoV01,
					PredicateType: slsa02.PredicateSLSAProvenance,
				},
				Predicate: slsa02.ProvenancePredicate{
					Builder: slsacommon.ProvenanceBuilder{
						ID: slsa.GithubHostedActionsBuilderID,
					},
					BuildType:   gh.TestBuildType,
					BuildConfig: gh.TestBuildConfig,
					Invocation: slsa02.ProvenanceInvocation{
						Environment: map[string]interface{}{
							"github_run_id":           "12345",
							"github_run_attempt":      "1",
							"github_actor":            "user",
							"github_base_ref":         "some/base_ref",
							"github_event_name":       "pull_request",
							"github_head_ref":         "some/head_ref",
							"github_ref":              "some/ref",
							"github_ref_type":         "branch",
							"github_repository_owner": "",
							"github_run_number":       "102937",
							"github_sha1":             "abcde",
						},
						ConfigSource: slsa02.ConfigSource{
							Digest: slsacommon.DigestSet{
								"sha1": "abcde",
							},
						},
					},
					Metadata: &slsa02.ProvenanceMetadata{
						BuildInvocationID: "12345-1",
					},
				},
			}
	g := slsa.NewHostedActionsGenerator(b).WithClients(&slsa.NilClientProvider{})
	if p, err := g.Generate(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else {
		if want, got := expected, p; !cmp.Equal(want, got) {
			t.Errorf("unexpected result\nwant: %#v\ngot:  %#v\ndiff: %v", want, got, cmp.Diff(want, got))
		}
	}
}
