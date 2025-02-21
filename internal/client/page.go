package client

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ConfluentPageBodyStorage struct {
	Representation string `json:"representation,omitempty"`
	Value          string `json:"value,omitempty"`
}

type ConfluentPageBody struct {
	Storage        ConfluentPageBodyStorage `json:"storage,omitempty"`
	AtlasDocFormat struct{}                 `json:"atlasDocFormat,omitempty"`
}

type ConfluentPageLinks struct {
	Webui    string `json:"webui,omitempty"`
	Editui   string `json:"editui,omitempty"`
	Edituiv2 string `json:"edituiv2,omitempty"`
	Tinyui   string `json:"tinyui,omitempty"`
}

type ConfluencePage struct {
	Id          string    `json:"id,omitempty"`
	Status      string    `json:"status,omitempty"`
	Title       string    `json:"title,omitempty"`
	SpaceId     string    `json:"space,omitempty"`
	ParentId    *string   `json:"parentId,omitempty"`
	ParentType  *string   `json:"parentType,omitempty"`
	Position    int       `json:"position,omitempty"`
	AuthorId    string    `json:"authorId,omitempty"`
	OwnerId     string    `json:"ownerId,omitempty"`
	LastOwnerId string    `json:"lastOwnerId,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	Version     struct {
		CreatedAt time.Time `json:"createdAt,omitempty"`
		Message   string    `json:"message,omitempty"`
		Number    int       `json:"number,omitempty"`
		MinorEdit bool      `json:"minorEdit,omitempty"`
		AuthorId  string    `json:"authorId,omitempty"`
	} `json:"version,omitempty"`
	Body  ConfluentPageBody  `json:"body,omitempty"`
	Links ConfluentPageLinks `json:"_links,omitempty"`
}

type ConfluencePageResponse struct {
	Results []ConfluencePage `json:"results,omitempty"`
	Links   struct {
		Next string `json:"next,omitempty"`
		Base string `json:"base,omitempty"`
	} `json:"_links,omitempty"`
}

func (c *ConfluenceClient) getPage(limit int, nextCursor string) ([]ConfluencePage, string, error) {
	params := map[string]string{
		"limit":       fmt.Sprintf("%d", limit),
		"body-format": "storage",
	}
	if nextCursor != "" {
		params["cursor"] = nextCursor
	}
	var response ConfluencePageResponse
	resp, err := c.newRequest("/wiki/api/v2/pages", params, &response)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.Results, response.Links.Next, nil
}

func (c *ConfluenceClient) GetAllPageStream(ch chan<- ConfluencePage, errCh chan<- error) {
	limit := 25
	nextCursor := ""
	for {
		pages, next, err := c.getPage(limit, nextCursor)
		if err != nil {
			errCh <- err
			continue
		}

		pagesId := []string{}
		for _, page := range pages {
			pagesId = append(pagesId, page.Id)
		}
		fmt.Println(pagesId)
		if next != "" {
			u, err := url.ParseQuery(strings.Split(next, "?")[1])
			if err != nil {
				errCh <- err
				continue
			}
			nextCursor = u.Get("cursor")
		}

		for _, page := range pages {
			ch <- page
		}
		if next == "" {
			break
		}

		// Rate limit to 10 request per second to avoid getting blocked
		time.Sleep(100 * time.Millisecond)
	}
}
