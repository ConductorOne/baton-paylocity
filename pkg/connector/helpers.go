package connector

import (
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func makeNextPageStr(offset, limit, total int) string {
	nextPage := makeNextPage(offset, limit, total)

	return strconv.Itoa(nextPage)
}

func makeNextPage(offset, limit, total int) int {
	n := offset + limit
	if n >= total {
		return -1 // sentinel value
	}

	return n
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
