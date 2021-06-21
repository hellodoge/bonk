package paired

import "errors"

type Error struct {
	Primary   error
	Secondary error
}

func NewPairedError(primary, secondary error) error {
	if primary == nil {
		return secondary
	} else if secondary == nil {
		return primary
	}
	return Error{
		Primary:   primary,
		Secondary: secondary,
	}
}

func (e Error) Error() string {
	return e.Primary.Error() + "; " + e.Secondary.Error()
}

func (e Error) Is(target error) bool {
	return errors.Is(e.Primary, target) || errors.Is(e.Secondary, target)
}
