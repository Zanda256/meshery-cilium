package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	walker "github.com/layer5io/meshkit/utils/walker"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultVersion string
var DefaultURL string
var DefaultGenerationMethod string
var WorkloadPath string
var MeshModelPath string
var AllVersions []string
var CRDNames []string

var meshmodelmetadata = map[string]interface{}{
	"Primary Color":   "#6B91C7",
	"Secondary Color": "#9AB0CF",
	"Shape":           "circle",
	"Logo URL":        "https://github.com/cncf/artwork/blob/master/projects/cilium/icon/white/cilium_icon-white.svg?short_path=d2fbc08",
	"SVG_Color":       "<svg id=\"Layer_1\" data-name=\"Layer 1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 300 304.21698\"><defs><style>.cls-1{fill:#cbdd72;}.cls-2{fill:#98ca3f;}.cls-3{fill:#6389c6;}.cls-4{fill:#e8282b;}.cls-5{fill:#f8c519;}.cls-6{fill:#f07525;}.cls-7{fill:#8162aa;}.cls-8{fill:#373737;}</style></defs><path class=\"cls-1\" d=\"M40.53139,62.5952h44.7715l22.38575,38.83793L85.30289,140.27106H40.53139L18.14564,101.43313Z\"/><path class=\"cls-2\" d=\"M40.53139,162.3871h44.7715l22.38575,38.56822L85.30289,239.52354H40.53139L18.14564,200.95532Z\"/><path class=\"cls-3\" d=\"M127.91673,211.47393h44.7715L195.074,250.04215l-22.38575,38.56822h-44.7715L105.531,250.04215Z\"/><path class=\"cls-4\" d=\"M127.91673,111.682h44.7715L195.074,150.52l-22.38575,38.83792h-44.7715L105.531,150.52Z\"/><path class=\"cls-5\" d=\"M127.91673,12.42954h44.7715L195.074,50.99776,172.68823,89.566h-44.7715L105.531,50.99776Z\"/><path class=\"cls-6\" d=\"M214.6278,62.5952h45.58062l22.79032,38.83793-22.79032,38.83793H214.6278l-22.79031-38.83793Z\"/><path class=\"cls-7\" d=\"M214.6278,162.3871h45.58062l22.79032,38.56822-22.79032,38.56822H214.6278l-22.79031-38.56822Z\"/><path class=\"cls-8\" d=\"M176.67452,104.66962h-53.4863L96.36091,150.94987l26.82731,45.95983h53.4863l27.03011-45.97708Zm-6.99407,79.91228H130.48651l-19.93464-33.56515,19.83322-34.01932h39.29536L189.497,151.01675Z\"/><path class=\"cls-8\" d=\"M176.67452,203.92211h-53.4863L96.36091,250.16028l26.82731,46.00191h53.4863l27.03011-46.00191Zm-6.99407,79.88853H130.48651l-19.93464-33.56621,19.83322-34.0382h39.29536l19.81652,34.0382Z\"/><path class=\"cls-8\" d=\"M176.67452,5.41714h-53.4863l-26.82731,46.297,26.82731,45.94311h53.4863l27.03011-45.94311ZM169.68045,85.372H130.48651L110.55187,51.71411l19.83322-33.99495h39.29536l19.74909,33.995Z\"/><path class=\"cls-8\" d=\"M264.05986,154.29587h-53.503l-26.81058,46.2797L210.55683,246.536h53.503L291.09,200.57557Zm-7.07875,79.89986h-39.194L197.95391,200.643l19.83324-34.07433h39.194l19.83377,34.00691Z\"/><path class=\"cls-8\" d=\"M264.05986,55.04338h-53.503l-26.81058,46.53377,26.81058,46.24573h53.503L291.09,101.57715Zm-7.07875,80.39667h-39.194l-19.83324-33.76149,19.83324-34.21837h39.194l19.83377,34.21837Z\"/><path class=\"cls-8\" d=\"M89.4251,154.29587H36.20743L9.515,200.57557,36.20743,246.536H89.4251l26.89419-45.96038Zm-7.02642,79.89986H43.40216L23.66818,200.643l19.734-34.07433H82.39868l19.734,34.07433Z\"/><path class=\"cls-8\" d=\"M89.4251,55.04338H36.20743L9.515,101.57715l26.69244,46.24573H89.4251l26.89419-46.24573Zm-7.02642,80.39667H43.40216l-19.734-33.79493,19.734-34.21892H82.39868l19.734,34.21892Z\"/></svg>",
	"SVG_White":       "<svg id=\"Layer_1\" data-name=\"Layer 1\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 300 304.21698\"><defs><style>.cls-1{fill:#fff;}</style></defs><path class=\"cls-1\" d=\"M176.67452,104.66962h-53.4863L96.36091,150.94987l26.82731,45.95983h53.4863l27.03011-45.97708Zm-6.99407,79.91228H130.48651l-19.93464-33.56515,19.83322-34.01932h39.29536L189.497,151.01675Z\"/><path class=\"cls-1\" d=\"M176.67452,203.92211h-53.4863L96.36091,250.16028l26.82731,46.00191h53.4863l27.03011-46.00191Zm-6.99407,79.88853H130.48651l-19.93464-33.56621,19.83322-34.0382h39.29536l19.81652,34.0382Z\"/><path class=\"cls-1\" d=\"M176.67452,5.41714h-53.4863l-26.82731,46.297,26.82731,45.94311h53.4863l27.03011-45.94311ZM169.68045,85.372H130.48651L110.55187,51.71411l19.83322-33.99495h39.29536l19.74909,33.995Z\"/><path class=\"cls-1\" d=\"M264.05986,154.29587h-53.503l-26.81058,46.2797L210.55683,246.536h53.503L291.09,200.57557Zm-7.07875,79.89986h-39.194L197.95391,200.643l19.83324-34.07433h39.194l19.83377,34.00691Z\"/><path class=\"cls-1\" d=\"M264.05986,55.04338h-53.503l-26.81058,46.53377,26.81058,46.24573h53.503L291.09,101.57715Zm-7.07875,80.39667h-39.194l-19.83324-33.76149,19.83324-34.21837h39.194l19.83377,34.21837Z\"/><path class=\"cls-1\" d=\"M89.4251,154.29587H36.20743L9.515,200.57557,36.20743,246.536H89.4251l26.89419-45.96038Zm-7.02642,79.89986H43.40216L23.66818,200.643l19.734-34.07433H82.39868l19.734,34.07433Z\"/><path class=\"cls-1\" d=\"M89.4251,55.04338H36.20743L9.515,101.57715l26.69244,46.24573H89.4251l26.89419-46.24573Zm-7.02642,80.39667H43.40216l-19.734-33.79493,19.734-34.21892H82.39868l19.734,34.21892Z\"/></svg>",
}

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category:    "Orchestration & Management",
	SubCategory: "Service Mesh",
	Metadata:    meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_CILIUM_SERVICE_MESH)],
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}

func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("cilium", "cilium")
	if len(AllVersions) == 0 {
		return
	}
	DefaultVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests

	//Get all the crd names
	w := walker.NewGithub()
	err := w.Owner("cilium").
		Repo("cilium").
		Branch("master").
		Root("pkg/k8s/apis/cilium.io/client/crds/v2/**").
		RegisterFileInterceptor(func(gca walker.GithubContentAPI) error {
			if gca.Content != "" {
				CRDNames = append(CRDNames, gca.Name)
			}
			return nil
		}).Walk()
	if err != nil {
		fmt.Println("Could not find CRD names. Will fail component creation...", err.Error())
	}
	DefaultURL = "https://raw.githubusercontent.com/cilium/cilium/" + "master" + "/pkg/k8s/apis/cilium.io/client/crds/v2/"
}
