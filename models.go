package eirini

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
)

const (
	//Environment Variable Names
	EnvDownloadURL        = "DOWNLOAD_URL"
	EnvBuildpacks         = "BUILDPACKS"
	EnvDropletUploadURL   = "DROPLET_UPLOAD_URL"
	EnvAppID              = "APP_ID"
	EnvStagingGUID        = "STAGING_GUID"
	EnvCompletionCallback = "COMPLETION_CALLBACK"
	EnvEiriniAddress      = "EIRINI_ADDRESS"

	RegisteredRoutes = "routes"
	OriginalRequest  = "original_request"

	RecipeBuildPacksDir    = "/var/lib/buildpacks"
	RecipeBuildPacksName   = "recipe-buildpacks"
	RecipeWorkspaceDir     = "/recipe_workspace"
	RecipeWorkspaceName    = "recipe-workspace"
	RecipeOutputName       = "staging-output"
	RecipeOutputLocation   = "/out"
	RecipePacksBuilderPath = "/packs/builder"

	CCUploaderInternalURL = "cc-uploader.service.cf.internal"
	CCCertsMountPath      = "/etc/config/certs"
	CCCertsVolumeName     = "cc-certs-volume"
	CCAPICertName         = "cc-server-crt"
	CCAPIKeyName          = "cc-server-crt-key"
	CCUploaderCertName    = "cc-uploader-crt"
	CCUploaderKeyName     = "cc-uploader-crt-key"
	CCInternalCACertName  = "internal-ca-cert"
)

type Config struct {
	Properties Properties `yaml:"opi"`
}

type Properties struct {
	ServerCertPath string `yaml:"server_cert_path"`
	ServerKeyPath  string `yaml:"server_key_path"`
	ServerCAPath   string `yaml:"server_ca_path"`

	KubeNamespace     string `yaml:"kube_namespace"`
	NatsPassword      string `yaml:"nats_password"`
	NatsIP            string `yaml:"nats_ip"`
	CcUploaderIP      string `yaml:"cc_uploader_ip"`
	CcInternalAPI     string `yaml:"cc_internal_api"`
	CCCertsSecretName string `yaml:"cc_certs_secret_name"`
	RegistryAddress   string `yaml:"registry_address"`
	EiriniAddress     string `yaml:"eirini_address"`
	DownloaderImage   string `yaml:"downloader_image"`
	UploaderImage     string `yaml:"uploader_image"`
	ExecutorImage     string `yaml:"executor_image"`

	LoggregatorAddress  string `yaml:"loggregator_address"`
	LoggregatorCertPath string `yaml:"loggergator_cert_path"`
	LoggregatorKeyPath  string `yaml:"loggregator_key_path"`
	LoggregatorCAPath   string `yaml:"loggregator_ca_path"`

	CCCertPath string `yaml:"cc_cert_path"`
	CCKeyPath  string `yaml:"cc_key_path"`
	CCCAPath   string `yaml:"cc_ca_path"`

	KubeConfigPath string `yaml:"kube_config_path"`

	RootfsVersion string `yaml:"rootfs_version"`
}

//go:generate counterfeiter . Stager
type Stager interface {
	Stage(string, cf.StagingRequest) error
	CompleteStaging(*models.TaskCallbackResponse) error
}

type StagerConfig struct {
	EiriniAddress   string
	DownloaderImage string
	UploaderImage   string
	ExecutorImage   string
}

//go:generate counterfeiter . Extractor
type Extractor interface {
	Extract(src, targetDir string) error
}

//go:generate counterfeiter . Bifrost
type Bifrost interface {
	Transfer(ctx context.Context, request cf.DesireLRPRequest) error
	List(ctx context.Context) ([]*models.DesiredLRPSchedulingInfo, error)
	Update(ctx context.Context, update cf.UpdateDesiredLRPRequest) error
	Stop(ctx context.Context, identifier opi.LRPIdentifier) error
	StopInstance(ctx context.Context, identifier opi.LRPIdentifier, index uint) error
	GetApp(ctx context.Context, identifier opi.LRPIdentifier) *models.DesiredLRP
	GetInstances(ctx context.Context, identifier opi.LRPIdentifier) ([]*cf.Instance, error)
}

func GetInternalServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s", appName)
}

func GetInternalHeadlessServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s-headless", appName)
}
