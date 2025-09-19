package gt_test

import (
	"fmt"
	"testing"

	"github.com/mitranim/gt"
)

func TestParseNullUuid(t *testing.T) {
	eq(``, gt.ParseNullUuid(``).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce018c4bef898ba4f293946049`).String())
	eq(`ddf1bfce018c4bef898ba4f293946049`, gt.ParseNullUuid(`ddf1bfce-018c-4bef-898b-a4f293946049`).String())
}

// TODO: test versioning.
func TestRandomNullUuid(t *testing.T) {
	eq(false, gt.RandomNullUuid().IsZero())
	eq(false, gt.RandomNullUuid().IsZero())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
	neq(gt.RandomNullUuid(), gt.RandomNullUuid())
}

// See `TestUuid_common`.
func TestNullUuid(t *testing.T) {
	t.Run(`GoString`, func(t *testing.T) {
		eq("gt.NullUuid{}", fmt.Sprintf(`%#v`, gt.NullUuid{}))
		eq("gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)", fmt.Sprintf(`%#v`, gt.ParseNullUuid(`b85ae23dc3f4468995d688e1ee645501`)))
	})
}
