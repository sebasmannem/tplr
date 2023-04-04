# tplr

Tekton Pipe Line Runner, a tool to run a tekton pipeline.
You hand it the namespace and name of the pipeline.
And you set environment variables to values you want to pass ass arguments.

And tplr wil:
- read the definition of the pipeline from kubernetes
- it creates a pipeline run from the definition
- it loops through the arguments and checks if they are set as environment variables and sets them accordingly
- it creates the pipelinerun
tracking the pipeline run can be done with tkn (for now)
