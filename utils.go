package keyrock

import "fmt"

func (c client) getURL(suffix string) string {
	return fmt.Sprintf("%s%s",c.options.BaseURL, suffix)
}
