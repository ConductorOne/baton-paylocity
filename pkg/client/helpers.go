package client

import (
	"net/url"
	"strconv"
	"strings"
)

const ItemsPerPage = 100

type ReqOpt func(reqURL *url.URL)

type PageOptions struct {
	Limit     int
	PageToken string
}

func withLimitParam(limit int) ReqOpt {
	if limit <= 0 {
		limit = ItemsPerPage
	}
	return withQueryParam("limit", strconv.Itoa(limit))
}

func withTotalCount() ReqOpt {
	return withQueryParam("includeTotalCount", "true")
}

func withOffset(offset string) ReqOpt {
	return withQueryParam("offset", offset)
}

func withIncludes(includes ...string) ReqOpt {
	return withQueryParam("include", strings.Join(includes, ","))
}

func withNextToken(token string) ReqOpt { return withQueryParam("nextToken", token) }

func getNextPageToken(offset, limit, total int) string {
	if offset+limit < total {
		return strconv.Itoa(offset + limit)
	}
	return ""
}

func withQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		if value != "" {
			q := reqURL.Query()
			q.Set(key, value)
			reqURL.RawQuery = q.Encode()
		}
	}
}
