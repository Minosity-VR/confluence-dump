package client

import (
	"testing"
)

func TestConfluenceClient_GetPage(t *testing.T) {
	c := NewConfluenceClient("company.atlassian.net", "cookie")
	got, next, err := c.getPage(1, "")
	t.Log(err)
	t.Log(next)
	t.Log(got)
}
