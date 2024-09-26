package char_test

import (
	"errors"
	"fmt"
	. "hello-go/char"
	"reflect"
	"testing"
)

func TestCharLen(t *testing.T) {
	noVetErrorf := CharLen

	wrapped := errors.New("inner error")
	for _, test := range []struct {
		err        error
		wantText   string
		wantUnwrap error
		wantSplit  []error
	}{{
		//err:        fmt.Errorf("%w", errors.New("should not happen")),
		wantUnwrap: nil,
	}, {
		err:        fmt.Errorf("%w", wrapped),
		wantText:   "inner error",
		wantUnwrap: wrapped,
	}, {
		err:        fmt.Errorf("added context: %w", wrapped),
		wantText:   "added context: inner error",
		wantUnwrap: wrapped,
	}, {
		err:        fmt.Errorf("%w with added context", wrapped),
		wantText:   "inner error with added context",
		wantUnwrap: wrapped,
	}, {
		err:        fmt.Errorf("%s %w %v", "prefix", wrapped, "suffix"),
		wantText:   "prefix inner error suffix",
		wantUnwrap: wrapped,
	}, {
		err:        fmt.Errorf("%[2]s: %[1]w", wrapped, "positional verb"),
		wantText:   "positional verb: inner error",
		wantUnwrap: wrapped,
	}, {
		err:      fmt.Errorf("%v", wrapped),
		wantText: "inner error",
	}, {
		err:      fmt.Errorf("added context: %v", wrapped),
		wantText: "added context: inner error",
	}, {
		err:      fmt.Errorf("%v with added context", wrapped),
		wantText: "inner error with added context",
	}, {
		err:        fmt.Errorf("%w", nil),
		wantText:   "%!w(<nil>)",
		wantUnwrap: nil, // still nil
	}} {
		if noVetErrorf(test.wantText) <= 0 && nil != test.err {
			t.Errorf("Formatted error: %v", test.err)
		} else if noVetErrorf(test.wantText) <= 0 && nil == test.err {
			continue
		}
		if got, want := errors.Unwrap(test.err), test.wantUnwrap; got != want {
			t.Errorf("Formatted error: %v\nerrors.Unwrap() = %v, want %v", test.err, got, want)
		}
		if got, want := splitErr(test.err), test.wantSplit; !reflect.DeepEqual(got, want) {
			t.Errorf("Formatted error: %v\nUnwrap() []error = %v, want %v", test.err, got, want)
		}
		if got, want := test.err.Error(), test.wantText; got != want {
			t.Errorf("err.Error() = %q, want %q", got, want)
		}
	}
}

func splitErr(err error) []error {
	if e, ok := err.(interface{ Unwrap() []error }); ok {
		return e.Unwrap()
	}
	return nil
}
