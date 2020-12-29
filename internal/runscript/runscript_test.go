package runscript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunscriptNew(t *testing.T) {
	rscript, err := New("path")
	assert.Nil(t, err)
	assert.Equal(t, rscript.Path, "path")
}

func TestRunscripStdOut(t *testing.T) {
	rscript := Script{StandardOut: []byte{}}
	assert.Equal(t, rscript.StandardOutput(), "")
}

func TestRunscriptRun(t *testing.T) {
	rscript := Script{Path: "/asdf"}
	err := rscript.Run()
	assert.Nil(t, err)
}
