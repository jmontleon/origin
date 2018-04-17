package automation_service_broker

import (
	"github.com/golang/glog"

	"github.com/openshift/origin/pkg/oc/bootstrap"
	"github.com/openshift/origin/pkg/oc/bootstrap/clusteradd/componentinstall"
	"github.com/openshift/origin/pkg/oc/bootstrap/docker/dockerhelper"
)

const (
	asbNamespace = "automation-broker-apb"
)

type AutomationServiceBrokerComponentOptions struct {
	InstallContext componentinstall.Context
}

func (c *AutomationServiceBrokerComponentOptions) Name() string {
	return "automation-service-broker"
}

func (c *AutomationServiceBrokerComponentOptions) Install(dockerClient dockerhelper.Interface, logdir string) error {
	params := map[string]string{
		"NAMESPACE": asbNamespace,
	}
	glog.V(2).Infof("instantiating automation service broker template with parameters %v", params)

	component := componentinstall.Template{
		Name:            "automation-service-broker",
		Namespace:       asbNamespace,
		InstallTemplate: bootstrap.MustAsset("install/automationservicebroker/deploy-automation-broker-apb.yaml"),
	}

	return component.MakeReady(
		c.InstallContext.ClientImage(),
		c.InstallContext.BaseDir(),
		params).Install(dockerClient, logdir)
}
