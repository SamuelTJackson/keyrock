package keyrock

import "fmt"

func (c client) getURL(posix string) string {
	return fmt.Sprintf("%s%s",c.options.BaseURL, posix)
}
