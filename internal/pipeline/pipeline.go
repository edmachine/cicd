package pipeline

import (
	"cicd/internal/runscript"
	"io/ioutil"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Step struct {
	Name     string   `yaml:"name"`
	Scripts  []string `yaml:"scripts"`
	Parallel bool     `yaml:"parallel"`
}

type Pipeline struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// NewFromPath will create a pipeline from the passed in path
func NewFromPath(path string) (pipeline Pipeline, err error) {
	// check if the path exists
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		zap.S().Error("err", err)
		return pipeline, err
	}
	err = yaml.Unmarshal(yamlFile, &pipeline)
	if err != nil {
		zap.S().Error("err", err)
		return pipeline, err
	}

	// zap.S().Infow("parsed from path", "pipeline", pipeline)

	return pipeline, nil
}

// Start will start running through the steps of the pipeline one step at a time
func (p Pipeline) Start() (err error) {
	zap.S().Infow("Starting Pipeline", "name", p.Name)
	// we will run steps in order and serially
	for _, step := range p.Steps {
		// run this step, if it fails return the error and stop the running pipeline
		if step.Parallel {
			if err = p.runParallelStep(step); err != nil {
				return err
			}
		} else if err = p.runStep(step); err != nil {
			return err
		}
	}
	return nil
}

// runHelper is an internal method that takes in the stepScript path and the waitgroup, and channels
// then runs the script and decrements the WG counter
func runHelper(path string, wg *sync.WaitGroup, errChan chan<- error, exitChan chan<- struct{}) {
	defer wg.Done()
	rscript, err := runscript.New(path)
	if err != nil {
		errChan <- err
		return
	}

	if err = rscript.Run(); err != nil {
		// error occured running the script
		errChan <- err
		return
	}
}

// runParallelStep is an internal method that will parallize running the scripts in a given step
func (p Pipeline) runParallelStep(step Step) (err error) {
	//setup channels to do communication
	// error channel will communicate errors back to this process
	errChan := make(chan error)
	// exit channel will cancel the running channels until they are all completed
	exitChan := make(chan struct{})
	// wait group to wait on the number of processes being run
	var wg sync.WaitGroup

	// spin off go routines to do the work
	for _, stepScript := range step.Scripts {
		zap.S().Infow("running script", "script", stepScript)
		wg.Add(1)
		go runHelper(stepScript, &wg, errChan, exitChan)
	}

	// wait for all of the goroutines to return
	wg.Wait()

	// get any errors off of the errChan
	if len(errChan) > 0 {
		for err := range errChan {
			zap.S().Error("err", err)
		}
		return err
	}

	return nil
}

// runStep is an internal method that is used to run scripts in a given step in serial
func (p Pipeline) runStep(step Step) (err error) {
	// zap.S().Infow("running step", "step", step)

	for _, stepScript := range step.Scripts {
		zap.S().Infow("running script", "script", stepScript)
		// create a new runscript object
		rscript, err := runscript.New(stepScript)
		// run in serial
		if err = rscript.Run(); err != nil {
			zap.S().Errorw(
				"error received from rscript.Run",
				"ExitCode", rscript.ExitCode,
				"StdErr", rscript.StandardErr,
				"StdOut", rscript.StandardOutput(),
			)
			return err
		}
	}
	return nil
}
