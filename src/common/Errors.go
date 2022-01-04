package common

import "errors"

var (
	ERR_LOCK_ALREADY_ACQUIRED = errors.New("lock already acquired")

	ERR_NO_LOCAL_IP_FOUND = errors.New("failed to find local ip")
)
