package template

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCenter(t *testing.T) {
	require.Equal(t, "fffff", center(3, "#", "fffff"))
	require.Equal(t, "##fffff###", center(10, "#", "fffff"))
	require.Equal(t, "###fffff###", center(11, "#", "fffff"))
}
