package dynamo

import "strings"

func ExtractKeyValue(value string) string {
	parts := strings.SplitN(value, "#", 2)

	if len(parts) == 2 {
		return parts[1]
	}

	return ""
}
