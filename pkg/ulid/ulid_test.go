package ulid_test

import (
	"strings"
	"testing"

	oklogulid "github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apexsudo/common/pkg/ulid"
)

const testPrefix = "usr"

func TestNew_prefixedWithGivenPrefix(t *testing.T) {
	t.Parallel()

	result := ulid.New(testPrefix)

	assert.True(t, strings.HasPrefix(result, testPrefix+"_"))
}

func TestNew_suffixIsValidULID(t *testing.T) {
	t.Parallel()

	result := ulid.New(testPrefix)

	parts := strings.SplitN(result, "_", 2)
	require.Len(t, parts, 2)

	_, err := oklogulid.ParseStrict(parts[1])

	require.NoError(t, err)
}

func TestNew_returnsUniqueValues(t *testing.T) {
	t.Parallel()

	first := ulid.New(testPrefix)
	second := ulid.New(testPrefix)

	assert.NotEqual(t, first, second)
}
