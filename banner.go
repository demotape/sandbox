package sandbox

import (
	"fmt"
)

func Welcome(demoName string, author string) string {
	var msg string
	msg = fmt.Sprintf("boxes -d shell -p a1l2 <(banner %s; echo Demo of: %s\necho Created by: %s)", "Demotape", demoName, author)

	return msg
}
