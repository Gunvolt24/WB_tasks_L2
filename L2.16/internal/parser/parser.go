package parser

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinksAndResources возвращает все ссылки и ресурсы с HTML
func ExtractLinksAndResources(body []byte, baseURL string) (links []string, resources []string, err error) {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse HTML: %v", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link := resolveURL(attr.Val, baseURL)
						if link != "" {
							links = append(links, link)
						}
					}
				}
			case "img", "script", "link":
				for _, attr := range n.Attr {
					if (n.Data == "img" && attr.Key == "src") ||
						(n.Data == "script" && attr.Key == "src") ||
						(n.Data == "link" && attr.Key == "href") {
						res := resolveURL(attr.Val, baseURL)
						if res != "" {
							resources = append(resources, res)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, resources, nil
}

// resolveURL делает абсолютный URL из относительного
func resolveURL(href, base string) string {
	if strings.HasPrefix(href, "javascript:") || href == "" || href[0] == '#' {
		return ""
	}
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseU, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return baseU.ResolveReference(u).String()
}
