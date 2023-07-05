package gt

import (
	"fmt"
	"io"
)

var (
	errInvalidChar    = fmt.Errorf(`invalid character`)
	errFormatMismatch = fmt.Errorf(`format mismatch`)
	errLengthMismatch = fmt.Errorf(`length mismatch`)
	errTerNullBool    = fmt.Errorf(`[gt] can't convert ternary null to boolean`)
	errUnrecLength    = fmt.Errorf(`unrecognized length`)
	errDigitEof       = fmt.Errorf(`expected digit, got %w`, io.EOF)
	errEmptySegment   = fmt.Errorf(`[gt] unexpected empty URL segment`)
)

func errParse(ptr *error, src string, typ string) {
	if *ptr != nil {
		*ptr = fmt.Errorf(`[gt] failed to parse %q into %v: %w`, src, typ, *ptr)
	}
}

func errInvalidCharAt(src string, ind int) error {
	for _, char := range src[ind:] {
		return fmt.Errorf(`[gt] invalid character %q in position %v`, char, ind)
	}
	return errInvalidChar
}

func errJsonString(src []byte, typ any) error {
	return fmt.Errorf(`[gt] can't decode %q into %T: expected string`, src, typ)
}

func errScanType(tar, inp any) error {
	return fmt.Errorf(`[gt] unrecognized input for type %T: type %T, value %v`, tar, inp, inp)
}

func errInvalidSegment(val string) error {
	return fmt.Errorf(`[gt] unexpected invalid URL segment %q`, val)
}
