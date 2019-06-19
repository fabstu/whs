package helper

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCurrentPackage(t *testing.T) {
	wantPackage := "aduu.dev/tools/aduu/helper"
	assert.Equal(t, CurrentPackage(), wantPackage, "should get correct package")
}