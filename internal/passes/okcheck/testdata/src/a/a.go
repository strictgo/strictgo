package a

import "strings"

func lookup() (int, bool)              { return 0, false }
func three() (int, int, bool)          { return 0, 0, false }
func withError() (int, error)          { return 0, nil }
func boolFirst() (bool, int)           { return false, 0 }
func single() bool                     { return false }
func generic[T any](value T) (T, bool) { return value, true }

type boolean bool

type boolAlias = bool

func namedBool() (int, boolean)     { return 0, false }
func aliasedBool() (int, boolAlias) { return 0, false }

type source struct{}

func (source) Lookup() (int, bool) { return 0, false }

type lookupSource interface {
	Lookup() (int, bool)
}

var packageValue, _ = lookup() // want "do not discard the boolean result"

func ignored(receiver source, dynamic lookupSource, indirect func() (int, bool)) {
	value, _ := lookup()    // want "do not discard the boolean result"
	value, _ = lookup()     // want "do not discard the boolean result"
	_, _ = lookup()         // want "do not discard the boolean result"
	value, _ = lookup()     // want "do not discard the boolean result"
	declared, _ := lookup() // want "do not discard the boolean result"
	grouped, _ := lookup()  // want "do not discard the boolean result"

	if value, _ := lookup(); value > 0 { // want "do not discard the boolean result"
		_ = value
	}
	value, _ = receiver.Lookup()                        // want "do not discard the boolean result"
	value, _ = dynamic.Lookup()                         // want "do not discard the boolean result"
	value, _ = indirect()                               // want "do not discard the boolean result"
	value, _ = generic(1)                               // want "do not discard the boolean result"
	value, _ = generic[int](1)                          // want "do not discard the boolean result"
	value, _ = namedBool()                              // want "do not discard the boolean result"
	value, _ = aliasedBool()                            // want "do not discard the boolean result"
	value, _ = func() (int, bool) { return 0, false }() // want "do not discard the boolean result"
	_, _ = strings.CutPrefix("prefix", "pre")           // want "do not discard the boolean result"
	_, _, _ = value, declared, grouped
}

func accepted(values map[string]int, input any, channel <-chan int) {
	value, ok := lookup()
	value, ok = lookup()
	_, ok = lookup()
	declared, preserved := lookup()
	value, _ = withError()
	_, value = boolFirst()
	_, _, _ = three()
	_, _ = 1, single()
	value, _ = values["key"]
	value, _ = input.(int)
	value, _ = <-channel
	_ = single()
	lookup()
	_, _, _, _ = value, ok, declared, preserved
}
