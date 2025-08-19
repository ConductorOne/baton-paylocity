package connector

import (
	"github.com/conductorone/baton-paylocity/pkg/client"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func getPageOptions(pToken *pagination.Token, defaultLimit int) client.PageOptions {
	if pToken == nil {
		return client.PageOptions{
			Limit:     defaultLimit,
			PageToken: "",
		}
	}

	limit := defaultLimit
	if pToken.Size > 0 {
		limit = pToken.Size
	}

	return client.PageOptions{
		Limit:     limit,
		PageToken: pToken.Token,
	}
}
