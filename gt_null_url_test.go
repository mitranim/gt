package gt_test

import (
	"fmt"
	"testing"

	"github.com/mitranim/gt"
)

func TestNullUrl(t *testing.T) {
	// Delegates to `gt.Join`. We just need to check the basics.
	t.Run(`AddPath`, func(t *testing.T) {
		eq(
			gt.NullUrl{Path: `one/two/three/four/five/six`},
			gt.NullUrl{Path: `one/two/three`}.AddPath(`four`, `five`, `/six`),
		)
	})

	// Delegates to `gt.Join`. We just need to check the basics.
	t.Run(`WithPath`, func(t *testing.T) {
		eq(
			gt.NullUrl{Path: `one/two/three`},
			gt.NullUrl{Path: `four`}.WithPath(`one`, `two`, `/three`),
		)
	})

	t.Run(`GoString`, func(t *testing.T) {
		eq(`gt.NullUrl{}`, fmt.Sprintf(`%#v`, gt.NullUrl{}))
		eq("gt.ParseNullUrl(`one://two.three/four?five=six#seven`)", fmt.Sprintf(`%#v`, gt.ParseNullUrl(`one://two.three/four?five=six#seven`)))
	})
}

func TestNullUrl_AddPath(t *testing.T) {
	invalid := func(src gt.NullUrl, more ...string) {
		panics(t, `[gt] unexpected empty URL segment`, func() { src.AddPath(more...) })
	}

	invalid(gt.NullUrl{}, ``)
	invalid(gt.NullUrl{Host: `two.three`}, ``)
	invalid(gt.NullUrl{Scheme: `one`, Host: `two.three`}, ``)
	invalid(gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `four`}, ``)
	invalid(gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/four`}, ``)
	invalid(gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/four`}, ``, `five`)
	invalid(gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/four`}, `five`, ``)

	same := func(src gt.NullUrl) { eq(src.AddPath(), src) }

	same(gt.NullUrl{}.AddPath())
	same(gt.NullUrl{Host: `one.two`}.AddPath())
	same(gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath())

	eq(
		gt.NullUrl{}.AddPath(`/`),
		gt.NullUrl{Path: `/`},
	)

	eq(
		gt.NullUrl{}.AddPath(`one`),
		gt.NullUrl{Path: `one`},
	)

	eq(
		gt.NullUrl{}.AddPath(`/one`),
		gt.NullUrl{Path: `/one`},
	)

	eq(
		gt.NullUrl{}.AddPath(`one`, `two`),
		gt.NullUrl{Path: `one/two`},
	)

	eq(
		gt.NullUrl{}.AddPath(`one/two`, `three`),
		gt.NullUrl{Path: `one/two/three`},
	)

	eq(
		gt.NullUrl{}.AddPath(`/one`, `two`),
		gt.NullUrl{Path: `/one/two`},
	)

	eq(
		gt.NullUrl{}.AddPath(`/one/two`, `three`),
		gt.NullUrl{Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{}.AddPath(`/one/two`, `/three`),
		gt.NullUrl{Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`/`),
		gt.NullUrl{Host: `two.three`, Path: `/`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`one`),
		gt.NullUrl{Host: `two.three`, Path: `/one`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`/one`),
		gt.NullUrl{Host: `two.three`, Path: `/one`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`one`, `two`),
		gt.NullUrl{Host: `two.three`, Path: `/one/two`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`one/two`, `three`),
		gt.NullUrl{Host: `two.three`, Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`/one`, `two`),
		gt.NullUrl{Host: `two.three`, Path: `/one/two`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`/one/two`, `three`),
		gt.NullUrl{Host: `two.three`, Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Host: `two.three`}.AddPath(`/one/two`, `/three`),
		gt.NullUrl{Host: `two.three`, Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`/`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`one`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`/one`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`one`, `two`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one/two`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`one/two`, `three`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`/one`, `two`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one/two`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`/one/two`, `three`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one/two/three`},
	)

	eq(
		gt.NullUrl{Scheme: `one`, Host: `two.three`}.AddPath(`/one/two`, `/three`),
		gt.NullUrl{Scheme: `one`, Host: `two.three`, Path: `/one/two/three`},
	)
}

func TestJoin_invalid(t *testing.T) {
	test := func(src, msg string, vals []string) {
		t.Helper()
		panics(t, msg, func() { gt.Join(vals...) })
	}

	test(``, `[gt] unexpected empty URL segment`, []string{``})
	test(``, `[gt] unexpected empty URL segment`, []string{`one`, ``})
	test(``, `[gt] unexpected empty URL segment`, []string{``, `one`})

	test(``, `[gt] unexpected invalid URL segment ".."`, []string{`..`})
	test(``, `[gt] unexpected invalid URL segment "/.."`, []string{`/..`})
	test(``, `[gt] unexpected invalid URL segment "../one"`, []string{`../one`})
	test(``, `[gt] unexpected invalid URL segment "/../one"`, []string{`/../one`})

	test(``, `[gt] unexpected invalid URL segment ".."`, []string{`one`, `..`})
	test(``, `[gt] unexpected invalid URL segment "/.."`, []string{`one`, `/..`})
	test(``, `[gt] unexpected invalid URL segment "../one"`, []string{`one`, `../one`})
	test(``, `[gt] unexpected invalid URL segment "/../one"`, []string{`one`, `/../one`})

	test(``, `[gt] unexpected invalid URL segment ".."`, []string{`..`, `one`})
	test(``, `[gt] unexpected invalid URL segment "/.."`, []string{`/..`, `one`})
	test(``, `[gt] unexpected invalid URL segment "../one"`, []string{`../one`, `one`})
	test(``, `[gt] unexpected invalid URL segment "/../one"`, []string{`/../one`, `one`})
}

func TestJoin_valid(t *testing.T) {
	test := func(exp string, vals []string) {
		t.Helper()
		eq(exp, gt.Join(vals...))
	}

	test(``, []string{})
	test(`/`, []string{`/`})
	test(`.one`, []string{`.one`})
	test(`/.one`, []string{`/.one`})
	test(`one`, []string{`one`})
	test(`one`, []string{`one/`})
	test(`/one`, []string{`/one`})
	test(`/one`, []string{`/one/`})
	test(`/one`, []string{`/one`})

	for _, one := range pathSegmentsOne {
		for _, two := range pathSegmentsTwo {
			for _, three := range pathSegmentsThree {
				test(`/one/two/three`, []string{`/`, one, two, three})
			}
		}
	}
}
