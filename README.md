# CICD
This CICD project is a mock to what a CICD golang application may do in order to achieve a concurrent pipeline.

This is a VERY Naive approach to a CICD Pipeline with the requirement that currently the pipelines are just scripts being run in order, but now we need to make a CICD tool that can accomplish this and be maintanable.

# Data layout
A pipeline is going to be a yaml file, each pipeline is going to have multiple steps, each step is going to contain multiple scripts and the option to run the scripts in parallel or in serial.

For example:
```
name: golang-app-1
steps:
    - name: gitpull
      scripts:
          -  "/bin/echo /path/to/gitpull.sh"
    - name: unit test
      parallel: true
      scripts:
          - "/bin/echo /path/to/unittest-1.sh"
          - "/bin/echo /path/to/unittest-2.sh"
    - name: build
      scripts:
          - "/bin/echo /path/to/build.sh"
    - name: integration-test
      parallel: true
      scripts:
          - "/bin/echo /path/to/integrationtest-1.sh"
          - "/bin/echo /path/to/integrationtest-2.sh"
    - name: package-binary
      scripts:
          - "/bin/echo /path/to/package.sh"
    - name: deploy-binary-to-qa
      scripts:
          - "/bin/echo /path/to/deploy.sh"
```

This means that the pipeline doesn't care about the scripts its running, just that they do not fail.
If they do fail, we should stop the pipeline and return an error up to the top level to either message the user, post it somewhere, or get the information out in a smart way

# Outcome
* This should be well designed testable software not a script.
* There should be a test suite which tests at least one function, preferably one that has some reasonable complexity such as the main pipeline controller.
* There should be a README.md with a discussion of your approach, trade-offs that were made, issues that should be addressed, assumptions that were made, etc.

I believe this is a starting point to a testable and designable piece of software with MVP in mind vs full feature/implementation.
There are a few unit tests currently written alongside the codebase, and as for integration testing we can easily make a handful of _normal_ yaml files that can be run through and ensure the stdout is the same each run.


# Discussion points
```
The core pipeline controller and it's supporting software are the main focus of this exercise in terms of code and tests not the scripts it executes and what they do (you don't need to know AWS, Docker, or Kubernetes), however, there should be some discussion in the README.md about other concerns such as how you might address the complexity of varying script API for each step, how concurrency may be handled, or how optional steps may be handled.
```

Currently we run 1 pipeline on each run of the program, instead of this in the main.go file we can easily add an API listening on a socket to take in any yaml file and run it with giving updates on where it is in the pipeline.

So the core pipeline as it is implemented in golang is working by running `/bin/echo <script path>`, this is just to get to a point where the code compiles, is testable, and can be extended easily. After this point we would adjust the code to run the script, or break the scripts up in the yaml file to pass into the `exec.Command(...)` function call in the runscript package.

Personally I would make a change to that as we would want to be able to change from using one-off unmaintainable scripts to using golang code to do the steps. Converting the one-off scripts into golang plugins (replace the runscript with runplugin, each step in the pipeline would instead have multiple plugins instead of scripts, and each plugin has a type). The plugins would need to have a Run method along with the same methods that our runscript package contains. This would allow us to convert the one-off scripts to be `golang-unit-test-plugin` and just have that type in the yaml file and then we are no longer calling arbitrary commands on the command line with arbitrary output.

As of this naive implementation we don't care about the output of the runscript run commands, but instead focus on the organizing and running of the commands in either parallel or serially. Optional scripts are not considered at this time. I think to add that we would want to figure out how to qualify the optional-ality (is this a word? hah) of the step. Why is it optional? Should we not break the execution if something failed vs making running the script optional? If a script fails we should have the option to running a `we failed, lets clean up` script.

The testing done also would need to be revamped when we remove the `/bin/echo` command, we would need to add a mock to catch the `runscript.Run()` method calls to exec.Command(...) and then be mocked out to return mocked data. This could also be a chance to create a mock library for the runscript package allowing anything relying on that package to be easily mockable.

Given the time constraints for implementing this I chose not to do that at this time, but it could be done fairly simply in the future to ensure we are loosely coupled to the runscript library in our other testing files.

The data layout I chose (pipeline contains multiple steps, each step contains multiple scripts) is a standard layout. This can also be enhanced and changed without much pain. Deciding on what the actual goals of this project require (is a UI going to show up, will we just run this via a webhook - how does this get triggered, how are we going to get this data back to a user of our application, etc ect) will shape what these data structures will look like moving forward, and will likely require a discussion internally and with the end users to determine how to build this product out.

As we continue to make enhancements to this, will we utilize a dynamic pipeline (not hard coded in yaml, but instead linked via an API and can change at a moments notice) or will we be able to make the assumptions that the pipelines dont change often enough to be something other than a list of plugins/order of operations that an end user would run. 


# End Points
I think the end goal should be to move away from the one-off scripts, build out plugins for a team that will use this, support an API to set the yaml pipeline file with a few endpoints to get the runs that have been run data (historical data in a database/cache), and then begin building the `plugins` that are needed by the teams who want to use this. Focus would be on an MVP of 1 team using the pipeline to run all of its testing, packaging, QA testing, and then deployment steps with the plugins that we need to build. After those plugins are completed, then we can move onto supporting more teams with different plugins (java vs golang unit testing plugin, etc).

This will scale fairly well until we hit a point where too many pipelines are running at once and we overload the system. We would then need to break out the pipelines API out from the run-the-pipeline machines. Once at that point it should be fine scaling up/out.
