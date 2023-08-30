package gt_test

import (
	"fmt"
	"testing"

	"github.com/mitranim/gt"
)

func TestTer(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq(`gt.TerNull`, fmt.Sprintf(`%#v`, gt.Ter(0)))
		eq(`gt.TerNull`, fmt.Sprintf(`%#v`, gt.TerNull))
		eq(`gt.TerFalse`, fmt.Sprintf(`%#v`, gt.TerFalse))
		eq(`gt.TerTrue`, fmt.Sprintf(`%#v`, gt.TerTrue))
		eq(`gt.Ter(3)`, fmt.Sprintf(`%#v`, gt.Ter(3)))
		eq(`gt.Ter(255)`, fmt.Sprintf(`%#v`, gt.Ter(255)))
	})
}
