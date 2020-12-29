package main

import (
	"cicd/internal/pipeline"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// feature flags for verbose logging and adding transaction names to the REPL prompt
	var (
		pipelinePath string
	)
	flag.StringVar(&pipelinePath, "pipeline", "test-pipeline.yaml", "the pipeline that you want to run through right now")
	flag.Parse()

	// set the global logger to be a production logger
	logger, _ := zap.NewProduction()
	_ = zap.ReplaceGlobals(logger)

	zap.S().Infow("Starting cicd server, with settings", "pipelinePath", pipelinePath)

	// start a prometheus metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println(http.ListenAndServe(":2112", nil))
	}()

	// TODO: make this an API endpoint to POST a yaml file instead of reading a single
	// pipeline from a file
	// create the pipeline from the passed in variable
	pipeline, err := pipeline.NewFromPath(pipelinePath)
	if err != nil {
		zap.S().Error("err", err)
		os.Exit(1)
	}

	// start the pipeline
	err = pipeline.Start()
	if err != nil {
		zap.S().Error("err", err)
		os.Exit(1)
	}
}
