package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineNew(t *testing.T) {
	pipeline, err := NewFromPath("testdata/test-pipeline.yaml")
	assert.Nil(t, err)
	assert.Equal(t, pipeline.Name, "golang-app-2")
	assert.Equal(t, len(pipeline.Steps), 1)
}

func TestPipelineStart(t *testing.T) {
	pipeline, err := NewFromPath("testdata/test-pipeline.yaml")
	assert.Nil(t, err)
	assert.Equal(t, pipeline.Name, "golang-app-2")
	assert.Equal(t, len(pipeline.Steps), 1)

	err = pipeline.Start()
	assert.Nil(t, err)
}

func TestPipelineParallel(t *testing.T) {
	pipeline, err := NewFromPath("testdata/test-pipeline-parallel.yaml")
	assert.Nil(t, err)
	assert.Equal(t, pipeline.Name, "golang-app-2")
	assert.Equal(t, len(pipeline.Steps), 1)
	assert.Equal(t, len(pipeline.Steps[0].Scripts), 3)

	err = pipeline.Start()
	assert.Nil(t, err)
}
