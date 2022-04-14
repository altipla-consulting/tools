package pr

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
	"libs.altipla.consulting/errors"

	"tools.altipla.consulting/cmd/ci/internal/login"
	"tools.altipla.consulting/cmd/ci/internal/query"
)

var client *github.Client

func initClient(ctx context.Context) error {
	if client != nil {
		return nil
	}

	auth, err := login.ReadAuthConfig()
	if err != nil {
		return errors.Trace(err)
	}
	if auth == nil {
		return errors.Errorf("Inicia sesión con `ci login` antes de interactuar con GitHub")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth.AccessToken})
	client = github.NewClient(oauth2.NewClient(ctx, ts))

	return nil
}

func ListBranches(ctx context.Context) ([]string, error) {
	if err := initClient(ctx); err != nil {
		return nil, errors.Trace(err)
	}

	org, err := query.CurrentOrg()
	if err != nil {
		return nil, errors.Trace(err)
	}
	repo, err := query.CurrentRepo()
	if err != nil {
		return nil, errors.Trace(err)
	}
	prs, _, err := client.PullRequests.List(ctx, org, repo, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var names []string
	for _, pr := range prs {
		names = append(names, pr.GetHead().GetRef())
	}
	return names, nil
}

func Create(ctx context.Context, title string) (string, error) {
	if err := initClient(ctx); err != nil {
		return "", errors.Trace(err)
	}

	org, err := query.CurrentOrg()
	if err != nil {
		return "", errors.Trace(err)
	}
	repo, err := query.CurrentRepo()
	if err != nil {
		return "", errors.Trace(err)
	}
	branch, err := query.CurrentBranch()
	if err != nil {
		return "", errors.Trace(err)
	}
	base, err := query.MainBranch()
	if err != nil {
		return "", errors.Trace(err)
	}

	req := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branch),
		Base:  github.String(base),
		Body:  github.String(""),
	}
	pr, _, err := client.PullRequests.Create(ctx, org, repo, req)
	if err != nil {
		return "", errors.Trace(err)
	}

	return pr.GetLinks().GetHTML().GetHRef(), nil
}

func Branch(ctx context.Context, id int64) (string, error) {
	if err := initClient(ctx); err != nil {
		return "", errors.Trace(err)
	}

	org, err := query.CurrentOrg()
	if err != nil {
		return "", errors.Trace(err)
	}
	repo, err := query.CurrentRepo()
	if err != nil {
		return "", errors.Trace(err)
	}
	pr, _, err := client.PullRequests.Get(ctx, org, repo, int(id))
	if err != nil {
		return "", errors.Trace(err)
	}

	return pr.GetHead().GetRef(), nil
}
