package core

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorArray []error

func (err ErrorArray) Error() string {
	var b bytes.Buffer

	for _, errs := range err {
		b.WriteString(fmt.Sprintf("%s, ", errs.Error()))
	}

	errs := b.String()
	return strings.TrimSuffix(errs, ", ")
}
