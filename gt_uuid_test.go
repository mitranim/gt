package gt_test

import (
	"fmt"
	"testing"

	"github.com/mitranim/gt"
)

func TestUuid(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq("gt.ParseUuid(`00000000000000000000000000000000`)", fmt.Sprintf(`%#v`, gt.Uuid{}))
		eq("gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseUuid(`b85ae23dc3f4468995d688e1ee645501`)))
	})
}
