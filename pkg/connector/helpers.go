package connector

import (
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func makeNextPage(offset, limit, total int) string {
	n := offset + limit
	if n >= total {
		return ""
	}
	return strconv.Itoa(n)
}

// parsePaginationToken - takes as pagination token and returns offset and limit in that order.
func parsePaginationToken(pToken *pagination.Token, defaultLimit int) (int, int, error) {
	limit := defaultLimit
	offset := 0

	if pToken != nil {
		if pToken.Size > 0 {
			limit = pToken.Size
		}

		if pToken.Token != "" {
			parsedOffset, err := strconv.Atoi(pToken.Token)
			if err != nil {
				return 0, 0, err
			}
			offset = parsedOffset
		}
	}
	return offset, limit, nil
}
