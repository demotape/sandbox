package sandbox

import (
	"fmt"
)

func Welcome(name string) string {
	var msg string
	msg = fmt.Sprintf("boxes -d shell -p a1l2 <(banner %s; echo Demo of: %s\necho Uploaded by: %s)", "Demotape", name, "Kim, Hirokuni")

	return msg
}
