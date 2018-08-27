// Package github accesses the GitHub API
package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var ghToken, _ = os.LookupEnv("CCOLLAGE_GH_TOKEN")

type Contributor struct {
	Username string `json:"login"`
	Avatar   string `json:"avatar_url"`
	URL      string `json:"html_url"`
}

func accessGitHub(url string, r *http.Request) ([]byte, http.Header, error) {
	ctx := appengine.NewContext(r)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; CCollage/0.0.1; +https://github.com/valentine/ccollage/)")

	// if ghToken != "" {
	// 	req.Header.Set("Authorization", fmt.Sprintf("token %v", ghToken))
	// }

	client := urlfetch.Client(ctx)
	response, err := client.Do(req)
	if err != nil {
		// log.Printf("ERROR: \n%+v", err.Error())
		return nil, nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// log.Printf("ERROR: \n%+v", err.Error())
		return nil, nil, err
	}

	return responseData, response.Header, nil
}

// buildURL builds the URL to access the GitHub API
func buildContributorURL(owner string, repo string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)
	return url
}

// getContributors retrieves a list of contributor usernames from the repository link
func getContributors(url string, r *http.Request) ([]Contributor, error) {
	resp, headers, err := accessGitHub(url, r)
	if err != nil {
		return nil, fmt.Errorf("Error:<br /><code>%+v</code>", err.Error())
	}

	var contributors []Contributor

	if json.Valid(resp) == true {
		err := json.Unmarshal(resp, &contributors)
		if err != nil {
			// log.Printf("ERROR: Unable to unmarshal JSON:\n%+v", string(resp))
			return nil, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(resp))
		}
	} else {
		// log.Printf("ERROR: GitHub JSON is not valid:\n%+v", string(resp))
		return nil, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(resp))
	}

	if headers.Get("Link") != "" {
		p := processHeaderLinkPages(headers.Get("Link"))
		for i := 2; i < p; i++ {
			pageURL := fmt.Sprintf("%s?page=%d", url, i)

			resp, _, err := accessGitHub(pageURL, r)
			if err != nil {
				return nil, fmt.Errorf("Error:<br /><code>%+v</code>", err.Error())
			}

			var c []Contributor

			if json.Valid(resp) == true {
				err := json.Unmarshal(resp, &c)
				if err != nil {
					// log.Printf("ERROR: Unable to unmarshal JSON:\n%+v", string(resp))
					return nil, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(resp))
				}
			} else {
				// log.Printf("ERROR: GitHub JSON is not valid:\n%+v", string(resp))
				return nil, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(resp))
			}

			contributors = append(contributors, c...)
		}
	}

	return contributors, nil
}

// GetAllContributors takes a repo and its owner and gets the relevant user information except Full Name
func GetAllContributors(owner string, repo string, r *http.Request) ([]Contributor, error) {
	url := buildContributorURL(owner, repo)
	c, err := getContributors(url, r)
	if err != nil {
		return c, err
	}
	return c, nil
}

func processHeaderLinkPages(headerLinks string) int {
	linkArray := strings.Split(headerLinks, ",")                // <url>; rel="next", <url>; rel="last"
	uLink := strings.Split(linkArray[1], ";")[0]                // <url>; rel="last"
	link := strings.Split(strings.Split(uLink, ">")[0], "<")[1] // <url>
	pages := strings.Split(link, "?page=")[1]                   // https://api.github.com/repositories/123/contributors?page=1
	num, err := strconv.Atoi(pages)
	if err != nil {
		log.Printf("ERROR: Unable to convert string to int: %+v", pages)
	}
	return num
}

func GetRateLimit(r *http.Request) (total int, remaining int, resettime int, err error) {
	url := "https://api.github.com/rate_limit"

	type RateLimit struct {
		Resources struct {
			Core struct {
				Total     int `json:"limit"`
				Remaining int `json:"remaining"`
				ResetTime int `json:"reset"`
			} `json:"core"`
		} `json:"resources"`
	}

	ctx := appengine.NewContext(r)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; CCollage/0.0.1; +https://github.com/valentine/ccollage/)")

	// if ghToken != "" {
	// 	req.Header.Set("Authorization", fmt.Sprintf("token %v", ghToken))
	// }

	client := urlfetch.Client(ctx)
	response, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", err.Error())
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", err.Error())
	}

	var rl RateLimit

	if json.Valid(responseData) == true {
		err := json.Unmarshal(responseData, &rl)
		if err != nil {
			// log.Printf("ERROR: Unable to unmarshal JSON:\n%+v", string(responseData))
			return 0, 0, 0, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(responseData))
		}
	} else {
		// log.Printf("ERROR: GitHub JSON is not valid:\n%+v", string(responseData))
		return 0, 0, 0, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(responseData))
	}

	total = rl.Resources.Core.Total
	remaining = rl.Resources.Core.Remaining
	resettime = rl.Resources.Core.ResetTime

	if total == 0 && resettime == 0 {
		return 0, 0, 0, fmt.Errorf("GitHub sent a bad response:<br /><code>%+v</code>", string(responseData))
	}

	return total, remaining, resettime, nil
}
