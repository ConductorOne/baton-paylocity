package connector

import (
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func getPageToken(pToken *pagination.Token) (string, error) {
	if pToken == nil {
		return "", nil
	}
	return pToken.Token, nil
}
