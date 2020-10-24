package syntask

import "errors"

var ErrInternalEmptyTask = errors.New("pipeline is finished, but task is null")
