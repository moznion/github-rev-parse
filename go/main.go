package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type authRoundTripper struct {
	rt    http.RoundTripper
	token string
}

func (art *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "token "+art.token)
	return art.rt.RoundTrip(req)
}

type shaResponse struct {
	SHA string `json:"sha"`
}

type refsResponse struct {
	Object shaResponse `json:"object"`
}

type apiClient struct {
	httpClient *http.Client
	baseURL    string
}

func (c *apiClient) getBranchSHA(key string) (string, error) {
	return c.getSHAViaRefs("heads", key)
}

func (c *apiClient) getTagSHA(key string) (string, error) {
	return c.getSHAViaRefs("tags", key)
}

func (c *apiClient) getSHAViaRefs(kind string, key string) (string, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/git/refs/%s/%s", c.baseURL, kind, key))
	if err == nil {
		defer resp.Body.Close()
		if 200 <= resp.StatusCode && resp.StatusCode <= 299 {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			var refs refsResponse
			if err := json.Unmarshal(b, &refs); err != nil {
				return "", err
			}
			return refs.Object.SHA, nil
		}
	}
	return "", nil
}

func (c *apiClient) getCommitSHA(key string) (string, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/commits/%s", c.baseURL, key))
	if err == nil {
		defer resp.Body.Close()
		if 200 <= resp.StatusCode && resp.StatusCode <= 299 {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			var sha shaResponse
			if err := json.Unmarshal(b, &sha); err != nil {
				return "", err
			}
			return sha.SHA, nil
		}
	}
	return "", nil
}

func main() {
	var githubToken string
	flag.StringVar(&githubToken, "token", "", "token of GitHub")
	flag.Parse()

	if len(flag.Args()) < 3 {
		fmt.Fprintf(os.Stderr, `ERROR: parameter(s) lacked
[usage]
  $ github-rev-parse <org> <repo> <key (commit hash, branch, tag)>
  options:
    --token=github-token : pass the token of GitHub
`)
		os.Exit(1)
	}

	org := flag.Arg(0)
	repo := flag.Arg(1)
	key := flag.Arg(2)

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	if githubToken != "" {
		httpClient.Transport = &authRoundTripper{
			rt:    http.DefaultTransport,
			token: githubToken,
		}
	}

	apiClient := &apiClient{
		httpClient: httpClient,
		baseURL:    fmt.Sprintf("https://api.github.com/repos/%s/%s", org, repo),
	}

	sha, err := apiClient.getBranchSHA(key)
	if err != nil {
		panic(err)
	}
	if sha != "" {
		fmt.Printf("%s\n", sha)
		return
	}

	sha, err = apiClient.getTagSHA(key)
	if err != nil {
		panic(err)
	}
	if sha != "" {
		fmt.Printf("%s\n", sha)
		return
	}

	sha, err = apiClient.getCommitSHA(key)
	if err != nil {
		panic(err)
	}
	if sha != "" {
		fmt.Printf("%s\n", sha)
		return
	}

	// there is no result that is matched
	os.Exit(1)
}
