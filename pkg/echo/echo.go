package echo

import (
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
)

func EchoHandler(request bothandler.Request) string {
	i := strings.ToLower(request.Content)
	uwu := uwucheck(i)
	if uwu != "" {
		return uwu
	}

	r, ok := echos[i]
	if ok {
		return r
	}

	for _, v := range fragments {
		a := strings.ToLower(v.From)

		c := false
		if a == v.From { // All lower
			c = strings.Contains(i, v.From)
		} else {
			c = strings.Contains(request.Content, v.From)
		}

		if c {
			if len(v.ExtraGuards) > 0 {
				count := 0
				for _, g := range v.ExtraGuards {
					if strings.Contains(i, g) {
						count++
					}
				}
				if count > 1 {
					return v.To
				}
				continue
			}

			if len(i) <= len(v.From)*7 {
				return v.To
			}
		}
	}
	return ""
}
