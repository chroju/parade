package cmd

import (
	"fmt"
	"strings"

	"github.com/chroju/parade/ssmctl"
)

func queryParser(query string) (string, string, error) {
	option := ssmctl.DescribeOptionEquals
	if query == "" {
		return query, option, nil
	}
	if strings.HasSuffix(query, "*") {
		query = strings.TrimPrefix(strings.TrimSuffix(query, "*"), "/")
		if strings.HasPrefix(query, "*") {
			query = strings.TrimPrefix(query, "*")
			option = ssmctl.DescribeOptionContains
		} else if !strings.Contains(query, "*") {
			option = ssmctl.DescribeOptionBeginsWith
		} else {
			return "", "", fmt.Errorf(ErrMsgQueryFormat)
		}
	}
	if strings.HasPrefix(query, "*") || strings.Contains(query, "*") {
		return "", "", fmt.Errorf(ErrMsgQueryFormat)
	}

	return query, option, nil
}
