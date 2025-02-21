package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type ConfluenceClient struct {
	BaseUrl *url.URL

	httpClient *http.Client
}

type Dumper struct {
	client *ConfluenceClient
}

func NewDumper(client *ConfluenceClient) *Dumper {
	return &Dumper{client: client}
}

func (d *Dumper) StartDumper(wg *sync.WaitGroup, fetchChan chan<- ConfluencePage, errChan chan<- error) {
	go func() {
		d.client.GetAllPageStream(fetchChan, errChan)
		wg.Done()
	}()
}

// Return a new ConfluenceClient
// confluenceHost: The host of the confluence instance. No need to pass the http/https scheme
// cookie: The cookie to use for authentication
func NewConfluenceClient(confluenceHost string, cookie string) *ConfluenceClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	jar.SetCookies(&url.URL{Scheme: "https", Host: confluenceHost}, []*http.Cookie{
		{
			Name:  "tenant.session.token",
			Value: cookie,
		},
	})
	return &ConfluenceClient{
		BaseUrl:    &url.URL{Scheme: "https", Host: confluenceHost},
		httpClient: &http.Client{Jar: jar},
	}
}

// Execute a query to the confluence host, using the set cookie
// path: The path to query
// params: The query parameters
// v: The object to unmarshal the response into
func (c *ConfluenceClient) newRequest(path string, params map[string]string, v interface{}) (*http.Response, error) {
	rel := &url.URL{Path: path}
	u := c.BaseUrl.ResolveReference(rel)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	u.RawQuery = p.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "confdump/0.0..1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
