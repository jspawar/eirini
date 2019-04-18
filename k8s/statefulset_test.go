package k8s_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/eirini"
	. "code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/k8s/k8sfakes"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
	"code.cloudfoundry.org/eirini/util/utilfakes"
	"code.cloudfoundry.org/eirini/rootfspatcher"
	"code.cloudfoundry.org/eirini/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	testcore "k8s.io/client-go/testing"
)

const (
	namespace               = "testing"
	timeout   time.Duration = 60 * time.Second
)

var _ = Describe("Statefulset", func() {

	var (
		err                   error
		client                *fake.Clientset
		statefulSetDesirer    opi.Desirer
		livenessProbeCreator  *k8sfakes.FakeProbeCreator
		readinessProbeCreator *k8sfakes.FakeProbeCreator
		hasher                *utilfakes.FakeHasher
		rootfsVersion         string
	)

	listStatefulSets := func() []appsv1.StatefulSet {
		list, listErr := client.AppsV1().StatefulSets(namespace).List(meta.ListOptions{})
		Expect(listErr).NotTo(HaveOccurred())
		return list.Items
	}

	getStatefulSet := func(lrp *opi.LRP) *appsv1.StatefulSet {
		labelSelector := fmt.Sprintf("guid=%s,version=%s", lrp.LRPIdentifier.GUID, lrp.LRPIdentifier.Version)
		ss, getErr := client.AppsV1().StatefulSets(namespace).List(meta.ListOptions{LabelSelector: labelSelector})
		Expect(getErr).NotTo(HaveOccurred())
		return &ss.Items[0]
	}

	BeforeEach(func() {
		client = fake.NewSimpleClientset()

		livenessProbeCreator = new(k8sfakes.FakeProbeCreator)
		readinessProbeCreator = new(k8sfakes.FakeProbeCreator)
		hasher = new(utilfakes.FakeHasher)
		hasher.HashReturns("random", nil)
		rootfsVersion = "version1"
	})

	JustBeforeEach(func() {
		statefulSetDesirer = &StatefulSetDesirer{
			Client:                client,
			Namespace:             namespace,
			RootfsVersion:         rootfsVersion,
			LivenessProbeCreator:  livenessProbeCreator.Spy,
			ReadinessProbeCreator: readinessProbeCreator.Spy,
			Hasher:                hasher,
		}
	})

	Context("When creating an LRP", func() {
		var lrp *opi.LRP

		Context("When app name only has [a-z0-9]", func() {
			JustBeforeEach(func() {
				livenessProbeCreator.Returns(&corev1.Probe{})
				readinessProbeCreator.Returns(&corev1.Probe{})
				lrp = createLRP("Baldur", "1234.5", "my.example.route")
				err = statefulSetDesirer.Desire(lrp)
			})

			It("should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should creates a healthcheck probe", func() {
				Expect(livenessProbeCreator.CallCount()).To(Equal(1))
			})

			It("should creates a readiness probe", func() {
				Expect(readinessProbeCreator.CallCount()).To(Equal(1))
			})

			It("should provide the process-guid to the pod annotations", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(statefulSet.Spec.Template.Annotations[cf.ProcessGUID]).To(Equal("Baldur-guid"))
			})

			It("should set name for the stateful set", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(statefulSet.Name).To(Equal("baldur-space-foo-random"))
			})

			It("should set space name as annotation on the statefulset", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(statefulSet.Annotations[cf.VcapSpaceName]).To(Equal("space-foo"))
			})

			It("should set podManagementPolicy to parallel", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(string(statefulSet.Spec.PodManagementPolicy)).To(Equal("Parallel"))
			})

			It("should set imagePullPolicy to Always", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(string(statefulSet.Spec.Template.Spec.Containers[0].ImagePullPolicy)).To(Equal("Always"))
			})

			It("should set rootfsVersion as a label", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(statefulSet.Labels).To(HaveKeyWithValue(rootfspatcher.RootfsVersionLabel, rootfsVersion))
				Expect(statefulSet.Spec.Template.Labels).To(HaveKeyWithValue(rootfspatcher.RootfsVersionLabel, rootfsVersion))
			})

			Context("When redeploying an existing LRP", func() {
				BeforeEach(func() {
					lrp = createLRP("Baldur", "1234.5", "my.example.route")
					_, createErr := client.AppsV1().StatefulSets(namespace).Create(toStatefulSet(lrp))
					Expect(createErr).ToNot(HaveOccurred())
				})

				It("should fail", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("When the app name contains unsupported characters", func() {
			JustBeforeEach(func() {
				lrp = createLRP("Балдър", "1234.5", "my.example.route")
				err = statefulSetDesirer.Desire(lrp)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should use the guid as a name", func() {
				statefulSet := getStatefulSet(lrp)
				Expect(statefulSet.Name).To(Equal("guid_1234-random"))
			})
		})
	})

	Context("When getting an app", func() {

		var (
			expectedLRP *opi.LRP
			actualLRP   *opi.LRP
		)

		BeforeEach(func() {
			expectedLRP = createLRP("Baldur", "1234.5", "my.example.route")
			// This is because toStatefulSet function mutates the metatdata map
			lrpToCreateStatefulSet := createLRP("Baldur", "1234.5", "my.example.route")
			_, createErr := client.AppsV1().StatefulSets(namespace).Create(toStatefulSet(lrpToCreateStatefulSet))
			Expect(createErr).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			actualLRP, err = statefulSetDesirer.Get(opi.LRPIdentifier{GUID: "guid_1234", Version: "version_1234"})
		})

		It("should not fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("return the expected LRP", func() {
			Expect(expectedLRP).To(Equal(actualLRP))
		})

		Context("when the app does not exist", func() {
			JustBeforeEach(func() {
				_, err = statefulSetDesirer.Get(opi.LRPIdentifier{GUID: "idontknow", Version: "42"})
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

		})
	})

	Context("When updating an app", func() {
		Context("when the app exists", func() {

			var (
				lrp *opi.LRP
			)

			BeforeEach(func() {

				lrp = createLRP("update", "7653.2", `["my.example.route"]`)

				statefulSet := toStatefulSet(lrp)
				_, createErr := client.AppsV1().StatefulSets(namespace).Create(statefulSet)
				Expect(createErr).NotTo(HaveOccurred())
			})

			Context("with replica count modified", func() {

				JustBeforeEach(func() {
					lrp.TargetInstances = 5
					lrp.Metadata = map[string]string{cf.LastUpdated: "123214.2"}
					err = statefulSetDesirer.Update(lrp)
				})

				It("scales the app without error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("updates the desired number of app instances", func() {
					Eventually(func() int32 {
						return *getStatefulSet(lrp).Spec.Replicas
					}).Should(Equal(int32(5)))
				})
			})

			Context("with modified routes", func() {

				JustBeforeEach(func() {
					lrp.Metadata = map[string]string{cf.VcapAppUris: `["my.example.route", "my.second.example.route"]`}
					err = statefulSetDesirer.Update(lrp)
				})

				It("should update the stored routes", func() {
					Eventually(func() string {
						return getStatefulSet(lrp).Annotations[eirini.RegisteredRoutes]
					}, 1*time.Second).Should(Equal(`["my.example.route", "my.second.example.route"]`))
				})
			})
		})

		Context("when the app does not exist", func() {

			JustBeforeEach(func() {
				err = statefulSetDesirer.Update(createLRP("name", "!234.0", "[something.strange]"))
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not create the app", func() {
				sets, listErr := client.AppsV1().StatefulSets(namespace).List(meta.ListOptions{})
				Expect(listErr).NotTo(HaveOccurred())
				Expect(sets.Items).To(BeEmpty())
			})
		})
	})

	Context("When listing apps", func() {

		var (
			actualLRPs   []*opi.LRP
			expectedLRPs []*opi.LRP
		)

		BeforeEach(func() {
			expectedLRPs = []*opi.LRP{
				createLRP("odin", "1234.5", "my.example.route"),
				createLRP("thor", "4567.8", "my.example.route"),
				createLRP("mimir", "9012.3", "my.example.route"),
			}
			// This is because toStatefulSet function mutates the metatdata map
			lrpsToCreateStatefulSets := []*opi.LRP{
				createLRP("odin", "1234.5", "my.example.route"),
				createLRP("thor", "4567.8", "my.example.route"),
				createLRP("mimir", "9012.3", "my.example.route"),
			}

			for _, l := range lrpsToCreateStatefulSets {
				statefulset := toStatefulSet(l)
				_, createErr := client.AppsV1().StatefulSets(namespace).Create(statefulset)
				Expect(createErr).ToNot(HaveOccurred())
			}
		})

		JustBeforeEach(func() {
			actualLRPs, err = statefulSetDesirer.List()
		})

		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("translates all existing statefulSets to opi.LRPs", func() {
			Expect(actualLRPs).To(ConsistOf(expectedLRPs))
		})

		Context("no statefulSets exist", func() {

			BeforeEach(func() {
				client = fake.NewSimpleClientset()
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an empy list of LRPs", func() {
				Expect(actualLRPs).To(BeEmpty())
			})
		})

		Context("fails to list the statefulsets", func() {

			BeforeEach(func() {
				reaction := func(action testcore.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("boom")
				}
				client.PrependReactor("list", "statefulsets", reaction)
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

		})
	})

	Context("Stop an LRP", func() {

		BeforeEach(func() {
			lrp := createLRP("Baldur", "1234.5", "my.example.route")
			_, err = client.AppsV1().StatefulSets(namespace).Create(toStatefulSet(lrp))
			Expect(err).ToNot(HaveOccurred())
		})

		It("deletes the statefulSet", func() {
			err = statefulSetDesirer.Stop(opi.LRPIdentifier{GUID: "guid_1234", Version: "version_1234"})
			Expect(err).ToNot(HaveOccurred())

			Eventually(listStatefulSets, timeout).Should(BeEmpty())
		})

		Context("when the statefulSet does not exist", func() {

			JustBeforeEach(func() {
				err = statefulSetDesirer.Stop(opi.LRPIdentifier{})
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Get LRP instances", func() {

		var (
			instances []*opi.Instance
			pod1      *corev1.Pod
			pod2      *corev1.Pod
		)

		BeforeEach(func() {
			since1 := meta.Unix(123, 0)
			pod1 = toPod("odin", 0, &since1)
			since2 := meta.Unix(456, 0)
			pod2 = toPod("odin", 1, &since2)
		})

		JustBeforeEach(func() {
			_, err = client.CoreV1().Pods(namespace).Create(pod1)
			Expect(err).ToNot(HaveOccurred())

			_, err = client.CoreV1().Pods(namespace).Create(pod2)
			Expect(err).ToNot(HaveOccurred())

			instances, err = statefulSetDesirer.GetInstances(opi.LRPIdentifier{GUID: "guid_1234", Version: "version_1234"})
		})

		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return the correct number of instances", func() {
			Expect(instances).To(HaveLen(2))
			Expect(instances[0]).To(Equal(toInstance(0, 123000000000, "RUNNING")))
			Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "RUNNING")))
		})

		Context("and time since creation is not available yet", func() {

			BeforeEach(func() {
				pod1 = toPod("mimir", 0, nil)
				since2 := meta.Unix(456, 0)
				pod2 = toPod("mimir", 1, &since2)
			})

			JustBeforeEach(func() {
				instances, err = statefulSetDesirer.GetInstances(opi.LRPIdentifier{GUID: "guid_1234", Version: "version_1234"})
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a default value", func() {
				Expect(instances).To(HaveLen(2))
				Expect(instances[0]).To(Equal(toInstance(0, 0, "RUNNING")))
				Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "RUNNING")))
			})
		})

		Context("and the pod has crashed", func() {
			BeforeEach(func() {
				pod1.Status.ContainerStatuses[0].State = corev1.ContainerState{
					Terminated: &corev1.ContainerStateTerminated{},
				}

				pod2.Status.ContainerStatuses[0].State = corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{},
				}
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a default value", func() {
				Expect(instances).To(HaveLen(2))
				Expect(instances[0]).To(Equal(toInstance(0, 123000000000, "CRASHED")))
				Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "CRASHED")))
			})
		})

		Context("and the StatefulSet was deleted/stopped", func() {

			BeforeEach(func() {
				event1 := &corev1.Event{
					Reason: "Killing",
					InvolvedObject: corev1.ObjectReference{
						Name:      "odin-0",
						Namespace: namespace,
						UID:       "odin-0-uid",
					},
				}
				event2 := &corev1.Event{
					Reason: "Killing",
					InvolvedObject: corev1.ObjectReference{
						Name:      "odin-1",
						Namespace: namespace,
						UID:       "odin-1-uid",
					},
				}

				event1.Name = "event1"
				event2.Name = "event2"

				_, clientErr := client.CoreV1().Events(namespace).Create(event1)
				Expect(clientErr).ToNot(HaveOccurred())
				_, clientErr = client.CoreV1().Events(namespace).Create(event2)
				Expect(clientErr).ToNot(HaveOccurred())
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a default value", func() {
				Expect(instances).To(HaveLen(0))
			})
		})

		Context("and the pod is pending", func() {
			BeforeEach(func() {
				pod1.Status.Phase = corev1.PodPending
				pod2.Status.ContainerStatuses[0].Ready = false
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a default value", func() {
				Expect(instances).To(HaveLen(2))
				Expect(instances[0]).To(Equal(toInstance(0, 123000000000, "CLAIMED")))
				Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "CLAIMED")))
			})
		})

		Context("and the pod phase is unknown", func() {
			BeforeEach(func() {
				pod1.Status.Phase = corev1.PodUnknown
				pod2.Status.Phase = corev1.PodUnknown
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a default value", func() {
				Expect(instances).To(HaveLen(2))
				Expect(instances[0]).To(Equal(toInstance(0, 123000000000, "UNKNOWN")))
				Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "UNKNOWN")))
			})
		})

		Context("the container status is not available yet", func() {

			BeforeEach(func() {
				pod1.Status.ContainerStatuses = []corev1.ContainerStatus{}
				pod2.Status.ContainerStatuses = []corev1.ContainerStatus{}
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an unknown status", func() {
				Expect(instances).To(HaveLen(2))
				Expect(instances[0]).To(Equal(toInstance(0, 123000000000, "UNKNOWN")))
				Expect(instances[1]).To(Equal(toInstance(1, 456000000000, "UNKNOWN")))
			})

		})

		Context("and the node has insufficient memory", func() {

			BeforeEach(func() {
				insufficientMemoryEvent := &corev1.Event{
					Reason:  "FailedScheduling",
					Message: "Some string including Insufficient memory",
					InvolvedObject: corev1.ObjectReference{
						Name:      "odin-0",
						Namespace: namespace,
						UID:       "odin-0-uid",
					},
				}

				_, clientErr := client.CoreV1().Events(namespace).Create(insufficientMemoryEvent)
				Expect(clientErr).ToNot(HaveOccurred())
			})

			It("shouldn't return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an unknown status", func() {
				Expect(instances).To(HaveLen(2))
				instance := toInstance(0, 123000000000, "UNCLAIMED")
				instance.PlacementError = "Insufficient resources: memory"
				Expect(instances).To(ContainElement(instance))
			})
		})
	})
})

func toPod(lrpName string, index int, time *meta.Time) *corev1.Pod {
	pod := corev1.Pod{}
	pod.Name = lrpName + "-" + strconv.Itoa(index)
	pod.UID = types.UID(pod.Name + "-uid")
	pod.Labels = map[string]string{
		"guid":    "guid_1234",
		"version": "version_1234",
	}

	pod.Status.StartTime = time
	pod.Status.Phase = corev1.PodRunning
	pod.Status.ContainerStatuses = []corev1.ContainerStatus{
		{
			State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
			Ready: true,
		},
	}
	return &pod
}

func toInstance(index int, since int64, state string) *opi.Instance {
	return &opi.Instance{
		Index: index,
		Since: since,
		State: state,
	}
}

func toStatefulSet(lrp *opi.LRP) *appsv1.StatefulSet {
	envs := MapToEnvVar(lrp.Env)
	fieldEnvs := []corev1.EnvVar{
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "CF_INSTANCE_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name: "CF_INSTANCE_INTERNAL_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}

	envs = append(envs, fieldEnvs...)
	ports := []corev1.ContainerPort{}
	for _, port := range lrp.Ports {
		ports = append(ports, corev1.ContainerPort{ContainerPort: port})
	}

	vols, volumeMounts := createVolumeSpecs(lrp.VolumeMounts)

	targetInstances := int32(lrp.TargetInstances)
	memory, err := resource.ParseQuantity(fmt.Sprintf("%dM", lrp.MemoryMB))
	if err != nil {
		panic(err)
	}
	cpu, err := resource.ParseQuantity(fmt.Sprintf("%dm", lrp.CPUWeight))
	if err != nil {
		panic(err)
	}

	automountServiceAccountToken := false

	namePrefix := fmt.Sprintf("%s-%s", lrp.AppName, lrp.SpaceName)
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name: fmt.Sprintf("%s-random", strings.ToLower(namePrefix)),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &targetInstances,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Annotations: map[string]string{
						cf.ProcessGUID: lrp.Metadata[cf.ProcessGUID],
						cf.VcapAppID:   lrp.Metadata[cf.VcapAppID],
					},
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: &automountServiceAccountToken,
					Containers: []corev1.Container{
						{
							Name:           "opi",
							Image:          lrp.Image,
							Command:        lrp.Command,
							Env:            envs,
							Ports:          ports,
							LivenessProbe:  &corev1.Probe{},
							ReadinessProbe: &corev1.Probe{},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: memory,
								},
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: memory,
									corev1.ResourceCPU:    cpu,
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: vols,
				},
			},
		},
	}

	statefulSet.Namespace = namespace

	labels := map[string]string{
		"guid":        lrp.GUID,
		"version":     lrp.Version,
		"source_type": "APP",
	}

	statefulSet.Spec.Template.Labels = labels

	statefulSet.Spec.Selector = &meta.LabelSelector{
		MatchLabels: labels,
	}

	statefulSet.Labels = labels

	statefulSet.Annotations = lrp.Metadata
	statefulSet.Annotations[eirini.RegisteredRoutes] = lrp.Metadata[cf.VcapAppUris]
	statefulSet.Annotations[cf.VcapSpaceName] = lrp.SpaceName

	return statefulSet
}

func createVolumeSpecs(lrpVolumeMounts []opi.VolumeMount) ([]corev1.Volume, []corev1.VolumeMount) {

	vols := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	for _, vol := range lrpVolumeMounts {
		vols = append(vols, corev1.Volume{
			Name: vol.ClaimName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: vol.ClaimName,
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      vol.ClaimName,
			MountPath: vol.MountPath,
		})
	}
	return vols, volumeMounts
}

func createLRP(name, lastUpdated, routes string) *opi.LRP {
	return &opi.LRP{
		LRPIdentifier: opi.LRPIdentifier{
			GUID:    "guid_1234",
			Version: "version_1234",
		},
		AppName:   name,
		SpaceName: "space-foo",
		Command: []string{
			"/bin/sh",
			"-c",
			"while true; do echo hello; sleep 10;done",
		},
		RunningInstances: 0,
		MemoryMB:         1024,
		Image:            "busybox",
		Ports:            []int32{8888, 9999},
		Metadata: map[string]string{
			cf.ProcessGUID: name + "-guid",
			cf.LastUpdated: lastUpdated,
			cf.VcapAppUris: routes,
			cf.VcapAppName: name,
			cf.VcapAppID:   "guid_1234",
			cf.VcapVersion: "version_1234",
		},
		VolumeMounts: []opi.VolumeMount{
			{
				ClaimName: "some-claim",
				MountPath: "/some/path",
			},
		},
	}
}
