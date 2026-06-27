// Package ulid wraps oklog/ulid to generate prefixed, time-sortable unique
// identifiers in the form "prefix_ULID" (e.g. "user_01JCMK8NVDBBZJCF7R2TWKHQJH").
package ulid

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

// New generates a ULID and returns it prefixed with the given label, producing
// strings like "user_01JCMK8NVDBBZJCF7R2TWKHQJH".  The prefix should be a
// short, stable identifier for the entity type (e.g. "user", "order").
func New(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, ulid.Make().String())
}
