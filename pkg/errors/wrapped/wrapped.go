package wrapped

import "errors"

type Error struct {
	Outer error
	Inner error
}

func (e Error) Error() string {
	return e.Outer.Error() + ": " + e.Inner.Error()
}

func (e Error) Is(target error) bool {
	if errors.Is(e.Outer, target) {
		return true
	} else if errors.Is(e.Inner, target) {
		return true
	} else {
		return false
	}
}