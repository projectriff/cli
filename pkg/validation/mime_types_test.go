package validation_test

import (
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	"github.com/projectriff/cli/pkg/validation"
)

func TestValidMimeType(t *testing.T) {
	expected := cli.EmptyFieldErrors
	actual := validation.MimeType("text/csv", "some-field")

	if diff := rifftesting.DiffFieldErrors(expected, actual); diff != "" {
		t.Errorf("(-expected, +actual): %s", diff)
	}
}

func TestInvalidMimeTypeWithMissingSlash(t *testing.T) {
	expected := cli.ErrInvalidValue("invalid", "some-field")
	actual := validation.MimeType("invalid", "some-field")

	if diff := rifftesting.DiffFieldErrors(expected, actual); diff != "" {
		t.Errorf("(-expected, +actual): %s", diff)
	}
}

func TestInvalidMimeTypeWithSingleTrailingSlash(t *testing.T) {
	expected := cli.ErrInvalidValue("invalid/", "some-field")
	actual := validation.MimeType("invalid/", "some-field")

	if diff := rifftesting.DiffFieldErrors(expected, actual); diff != "" {
		t.Errorf("(-expected, +actual): %s", diff)
	}
}
