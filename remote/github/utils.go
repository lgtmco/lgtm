package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func setupClient(rawurl, accessToken string) *github.Client {
	token := oauth2.Token{AccessToken: accessToken}
	source := oauth2.StaticTokenSource(&token)
	client := oauth2.NewClient(oauth2.NoContext, source)
	github := github.NewClient(client)
	github.BaseURL, _ = url.Parse(rawurl)
	return github
}

// GetHook is a heper function that retrieves a hook by
// hostname. To do this, it will retrieve a list of all hooks
// and iterate through the list.
func GetHook(client *github.Client, owner, name, rawurl string) (*github.Hook, error) {
	hooks, _, err := client.Repositories.ListHooks(owner, name, nil)
	if err != nil {
		return nil, err
	}
	newurl, err := url.Parse(rawurl)
	if err != nil {
		fmt.Println("error parsing new hook url", rawurl, err)
		return nil, err
	}
	for _, hook := range hooks {
		hookurl, ok := hook.Config["url"].(string)
		if !ok {
			continue
		}
		oldurl, err := url.Parse(hookurl)
		if err != nil {
			fmt.Println("error parsing old hook url", hookurl, err)
			continue
		}
		if newurl.Host == oldurl.Host {
			return &hook, nil
		}
	}
	return nil, nil
}

func DeleteHook(client *github.Client, owner, name, url string) error {
	hook, err := GetHook(client, owner, name, url)
	if err != nil {
		return err
	}
	if hook == nil {
		return nil
	}
	_, err = client.Repositories.DeleteHook(owner, name, *hook.ID)
	return err
}

// CreateHook is a heper function that creates a post-commit hook
// for the specified repository.
func CreateHook(client *github.Client, owner, name, url string) (*github.Hook, error) {
	var hook = new(github.Hook)
	hook.Name = github.String("web")
	hook.Events = []string{"issue_comment"}
	hook.Config = map[string]interface{}{}
	hook.Config["url"] = url
	hook.Config["content_type"] = "json"
	created, _, err := client.Repositories.CreateHook(owner, name, hook)
	return created, err
}

// GetFile is a heper function that retrieves a file from
// GitHub and returns its contents in byte array format.
func GetFile(client *github.Client, owner, name, path, ref string) ([]byte, error) {
	var opts = new(github.RepositoryContentGetOptions)
	opts.Ref = ref
	content, _, _, err := client.Repositories.GetContents(owner, name, path, opts)
	if err != nil {
		return nil, err
	}
	return content.Decode()
}

// GetUserRepos is a helper function that returns a list of
// all user repositories. Paginated results are aggregated into
// a single list.
func GetUserRepos(client *github.Client) ([]github.Repository, error) {
	var repos []github.Repository
	var opts = github.RepositoryListOptions{}
	opts.PerPage = 100
	opts.Page = 1

	// loop through user repository list
	for opts.Page > 0 {
		list, resp, err := client.Repositories.List("", &opts)
		if err != nil {
			return nil, err
		}
		repos = append(repos, list...)

		// increment the next page to retrieve
		opts.Page = resp.NextPage
	}

	return repos, nil
}
