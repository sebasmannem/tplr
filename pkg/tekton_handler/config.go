package tekton_handler

import (
	"github.com/tektoncd/cli/pkg/actions"
	"github.com/tektoncd/cli/pkg/cli"
	prcmd "github.com/tektoncd/cli/pkg/cmd/pipelinerun"
	"github.com/tektoncd/cli/pkg/options"
	"github.com/tektoncd/cli/pkg/pipelinerun"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"strings"
	"time"
)

var (
	pipelineGroupResource = schema.GroupVersionResource{Group: "tekton.dev", Resource: "pipelines"}
	//pipelineRunGroupResource = schema.GroupVersionResource{Group: "tekton.dev", Resource: "pipelineruns"}
)

type Config struct {
	Pipeline    string `yaml:"pipeline"`
	KubeConfig  string `yaml:"kube_config"`
	KubeContext string `yaml:"kube_context"`
	Namespace   string `yaml:"namespace"`
	Timeout     string `yaml:"pipeline_timeout"`
	clients     *cli.Clients
}

func (c *Config) getTektonParams() (*cli.TektonParams, error) {
	params := &cli.TektonParams{}
	params.SetKubeConfigPath(c.KubeConfig)
	params.SetKubeContext(c.KubeContext)
	params.SetNamespace(c.Namespace)
	return params, nil
}

func (c *Config) Connect() error {
	if c.clients != nil {
		log.Info("Already connected")
		return nil
	}
	if params, err := c.getTektonParams(); err != nil {
		return err
	} else if clients, err := params.Clients(); err != nil {
		return err
	} else {
		c.clients = clients
	}
	return nil
}

func (c *Config) GetPipeline() (*v1.Pipeline, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}
	var pipeline v1.Pipeline
	gvr, err := actions.GetGroupVersionResource(pipelineGroupResource, c.clients.Tekton.Discovery())
	if err != nil {
		return nil, err
	}

	if gvr.Version == "v1" {
		err := actions.GetV1(pipelineGroupResource, c.clients, c.Pipeline, c.Namespace, metav1.GetOptions{}, &pipeline)
		if err != nil {
			return nil, err
		}
		return &pipeline, nil

	}

	var pipelineV1beta1 v1beta1.Pipeline
	err = actions.GetV1(pipelineGroupResource, c.clients, c.Pipeline, c.Namespace, metav1.GetOptions{}, &pipelineV1beta1)
	if err != nil {
		return nil, err
	}
	err = pipelineV1beta1.ConvertTo(ctx, &pipeline)
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

func plParamsToPrParams(plParams v1beta1.ParamSpecs) v1beta1.Params {
	var prParams v1beta1.Params
	var paramValue string
	for _, plParam := range plParams {
		if envVal := os.Getenv(strings.ToUpper(plParam.Name)); envVal != "" {
			paramValue = envVal
		} else {
			paramValue = plParam.Default.StringVal
		}
		prParams = append(prParams, v1beta1.Param{
			Name: plParam.Name,
			Value: v1beta1.ParamValue{
				Type:      plParam.Default.Type,
				StringVal: paramValue,
				ArrayVal:  plParam.Default.ArrayVal,
				ObjectVal: plParam.Default.ObjectVal,
			},
		})
	}
	return prParams
}

func (c *Config) Run() error {
	var pl v1beta1.Pipeline
	if pipeline, err := c.GetPipeline(); err != nil {
		return err
	} else if err = pl.ConvertFrom(ctx, pipeline); err != nil {
		return err
	} else {
		pr := &v1beta1.PipelineRun{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "tekton.dev/v1beta1",
				Kind:       "PipelineRun",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: c.Namespace,
			},
			Spec: v1beta1.PipelineRunSpec{
				PipelineRef: &v1beta1.PipelineRef{Name: pipeline.ObjectMeta.Name},
			},
		}
		pr.ObjectMeta.GenerateName = pipeline.ObjectMeta.Name + "-run-"
		if c.Timeout != "" {
			timeoutDuration, err := time.ParseDuration(c.Timeout)
			if err != nil {
				return err
			}
			pr.Spec.Timeouts.Pipeline = &metav1.Duration{Duration: timeoutDuration}
		}
		pr.ObjectMeta.Labels = pipeline.Labels
		pr.Spec.Params = plParamsToPrParams(pl.Spec.Params)

		if prCreated, err := pipelinerun.Create(c.clients, pr, metav1.CreateOptions{}, c.Namespace); err != nil {
			return err
		} else if params, err := c.getTektonParams(); err != nil {
			return err
		} else {
			log.Infof("Created pipelinerun %s", prCreated.GetName())
			log.Infoln("Waiting for logs to be available...")
			runLogOpts := &options.LogOptions{
				PipelineName:    pipeline.ObjectMeta.Name,
				PipelineRunName: prCreated.Name,
				Stream: &cli.Stream{
					Out: os.Stdout,
					Err: os.Stderr,
				},
				Follow:    true,
				Prefixing: true,
				Params:    params,
				AllSteps:  false,
			}
			return prcmd.Run(runLogOpts)
		}
	}
}
