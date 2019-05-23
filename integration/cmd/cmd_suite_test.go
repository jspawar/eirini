package cmd_test

import (
	"io/ioutil"
	"os"
	"testing"

	"code.cloudfoundry.org/eirini"
	natsserver "github.com/nats-io/gnatsd/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cmd Suite")
}

func pathToTestFixture(relativePath string) string {
	cwd, err := os.Getwd()
	Expect(err).ToNot(HaveOccurred())
	return cwd + "/fixtures/" + relativePath
}

func defaultEiriniConfig(natsServerOpts natsserver.Options) *eirini.Config {
	config := &eirini.Config{
		Properties: eirini.Properties{
			KubeConfigPath:      pathToTestFixture("kube.conf"),
			CCCAPath:            pathToTestFixture("cert"),
			CCCertPath:          pathToTestFixture("cert"),
			CCKeyPath:           pathToTestFixture("key"),
			NatsIP:              natsServerOpts.Host,
			NatsPassword:        natsServerOpts.Password,
			LoggregatorCAPath:   pathToTestFixture("cert"),
			LoggregatorCertPath: pathToTestFixture("cert"),
			LoggregatorKeyPath:  pathToTestFixture("key"),
		},
	}

	return config
}

func createOpiConfigFromFixtures(config *eirini.Config) (*os.File, error) {
	configFile, err := ioutil.TempFile("", "opi-config.yml")
	Expect(err).ToNot(HaveOccurred())

	bs, err := yaml.Marshal(config)
	Expect(err).ToNot(HaveOccurred())
	err = ioutil.WriteFile(configFile.Name(), bs, os.ModePerm)
	Expect(err).ToNot(HaveOccurred())

	return configFile, err
}
