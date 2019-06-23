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

	resp, err := httpClient.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", org, repo, key))
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	if 200 <= resp.StatusCode && resp.StatusCode <= 299 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			os.Exit(1)
		}
		var sha shaResponse
		if err := json.Unmarshal(b, &sha); err != nil {
			os.Exit(1)
		}
		fmt.Printf("%s\n", sha.SHA)
		return
	}

	// there is no result that is matched
	os.Exit(1)
}
