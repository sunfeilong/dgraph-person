package tools

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestIsNumber(t *testing.T) {

    match1 := IsNumber("dd")
    assert.False(t, match1, "dd is not number")

    match2 := IsNumber("#")
    assert.False(t, match2, "# is not number")

}
