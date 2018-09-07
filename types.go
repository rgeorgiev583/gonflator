package gonflator

import (
	"fmt"
)

type ErrorCode int

const (
	CouldNotInitialize ErrorCode = -2
	NotFound                     = -1

	// No error
	NoError = 0

	// Out of memory
	ENOMEM = 1

	// Internal error
	EINTERNAL = 2

	// Invalid path expression
	EPATHX = 3

	// No match for path expression
	ENOMATCH = 4

	// Too many matches for path expression
	EMMATCH = 5

	// Cannot move node into its descendant
	EMVDESC = 10

	// Invalid argument in function call
	EBADARG = 12
)

type Error struct {
	Code ErrorCode

	// Human-readable error message
	Message string

	// Human-readable message elaborating the error. For example, when
	// the error code is EPATHX, this will explain how the path
	// expression is invalid
	MinorMessage string

	// Details about the error. For example, for EPATHX, indicates
	// where in the path expression the error occurred.
	Details string
}

func (err Error) Error() string {
	return fmt.Sprintf("Message: %s - Minor message: %s - Details: %s",
		err.Message, err.MinorMessage, err.Details)
}
