package hub

import (
	"errors"
)

var (
	ErrOptResendLimitInvalid = errors.New("invalid ResendLimit value")
	ErrOptLocalAddrInvalid   = errors.New("invalid LocalAddr value")
)
