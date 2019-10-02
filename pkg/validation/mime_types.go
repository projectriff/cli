package validation

import (
	"strings"

	"github.com/projectriff/cli/pkg/cli"
)

func MimeType(mimeType, field string) cli.FieldErrors {
	errs := cli.EmptyFieldErrors

	index := strings.Index(mimeType, "/")
	if index == -1 || index == len(mimeType)-1 {
		errs = errs.Also(cli.ErrInvalidValue(mimeType, field))
	}

	return errs
}
