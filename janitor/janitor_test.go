package janitor

import (
	"testing"
)

func TestClean(t *testing.T) {

	if Clean() != 1 {
		t.Error("error!!!")
	}
}
