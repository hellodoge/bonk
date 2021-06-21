package paired

import "errors"

type Error struct {
	Primary   error
	Secondary error
}

func (e Error) Error() string {
	return e.Primary.Error() + "; " + e.Secondary.Error()
}

func (e Error) Is(target error) bool {
	return errors.Is(e.Primary, target) || errors.Is(e.Secondary, target)
}
