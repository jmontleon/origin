package ansible_service_broker

import (
	"encoding/base64"
	"github.com/golang/glog"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/openshift/origin/pkg/cmd/util/variable"
	"github.com/openshift/origin/pkg/oc/bootstrap"
	"github.com/openshift/origin/pkg/oc/bootstrap/clusteradd/componentinstall"
	"github.com/openshift/origin/pkg/oc/bootstrap/clusterup/kubeapiserver"
	"github.com/openshift/origin/pkg/oc/bootstrap/docker/dockerhelper"
	"github.com/openshift/origin/pkg/oc/errors"
)

const (
	asbNamespace = "ansible-service-broker"
)

type AnsibleServiceBrokerComponentOptions struct {
	InstallContext componentinstall.Context
}

func (c *AnsibleServiceBrokerComponentOptions) Name() string {
	return "ansible-service-broker"
}

func (c *AnsibleServiceBrokerComponentOptions) Install(dockerClient dockerhelper.Interface, logdir string) error {
	imageTemplate := variable.NewDefaultImageTemplate()
	imageTemplate.Format = c.InstallContext.ImageFormat()
	imageTemplate.Latest = false

	masterConfigDir := path.Join(c.InstallContext.BaseDir(), kubeapiserver.KubeAPIServerDirName)
	serviceCABytes, err := ioutil.ReadFile(filepath.Join(masterConfigDir, "service-signer.crt"))
	serviceCAString := base64.StdEncoding.EncodeToString(serviceCABytes)
	if err != nil {
		return errors.NewError("unable to read service signer cert").WithCause(err)
	}

	params := map[string]string{
		"IMAGE":          imageTemplate.ExpandOrDie("ansible-service-broker"),
		"NAMESPACE":      asbNamespace,
		"BROKER_CA_CERT": serviceCAString,
	}
	glog.V(2).Infof("instantiating ansible service broker ansible with parameters %v", params)

	component := componentinstall.Template{
		Name:            "ansible-service-broker",
		Namespace:       asbNamespace,
		InstallTemplate: bootstrap.MustAsset("install/ansibleservicebroker/deploy-ansible-service-broker.template.yaml"),
	}

	err = component.MakeReady(
		c.InstallContext.ClientImage(),
		c.InstallContext.BaseDir(),
		params).Install(dockerClient, logdir)

	if err != nil {
		return err
	}

	return nil
}
