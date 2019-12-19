package release

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"

	shipper "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	"github.com/bookingcom/shipper/pkg/util/conditions"
)

const (
	clusterName   = "minikube"
	namespace     = "test-namespace"
	incumbentName = "0.0.1"
	contenderName = "0.0.2"
)

func init() {
	conditions.StrategyConditionsShouldDiscardTimestamps = true
}

var app = &shipper.Application{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test-app",
		Namespace: namespace,
		UID:       "foobarbaz",
	},
	Status: shipper.ApplicationStatus{
		History: []string{incumbentName, contenderName},
	},
}

// TestCompleteStrategyNoController tests the complete "vanguard" strategy,
// end-to-end.  This test exercises only the Executor.execute() method, using
// hard coded incumbent and contender releases, checking if the generated
// patches were the expected for the strategy at a given moment. Since the
// target objects are always in a Ready state, there's no need to mimic the
// other controllers, as just updating their specs will trigger the necessary
// conditions.
func TestCompleteStrategyNoController(t *testing.T) {
	totalReplicaCount := uint(10)
	executor := &Executor{
		contender: buildContender(totalReplicaCount),
		incumbent: buildIncumbent(totalReplicaCount),
		recorder:  record.NewFakeRecorder(42),
		strategy:  vanguard,
	}

	for _, step := range []int32{1, 2} {
		executor.contender.release.Spec.TargetStep = step

		// Execute first part of strategy's step.
		if newSpec, err := ensureCapacityPatch(executor, contenderName, Contender); err != nil {
			t.Fatal(err)
		} else {
			executor.contender.capacityTarget.Spec = *newSpec
		}

		// Execute second part of strategy's step.
		if newSpec, err := ensureTrafficPatch(executor, contenderName, Contender); err != nil {
			t.Fatal(err)
		} else {
			executor.contender.trafficTarget.Spec = *newSpec
		}

		// Execute third part of strategy's step.
		if newSpec, err := ensureTrafficPatch(executor, incumbentName, Incumbent); err != nil {
			t.Fatal(err)
		} else {
			executor.incumbent.trafficTarget.Spec = *newSpec
		}

		// Execute fourth part of strategy's step.
		if newSpec, err := ensureCapacityPatch(executor, incumbentName, Incumbent); err != nil {
			t.Fatal(err)
		} else {
			executor.incumbent.capacityTarget.Spec = *newSpec
		}

		// Execute fifth part of strategy's step.
		if newStatus, err := ensureReleasePatch(executor, contenderName); err != nil {
			t.Fatal(err)
		} else {
			executor.contender.release.Status = *newStatus
		}
	}

	if err := ensureFinalReleasePatches(executor); err != nil {
		t.Fatal(err)
	}

}

// buildIncumbent returns a releaseInfo with an incumbent release and
// associated objects.
func buildIncumbent(totalReplicaCount uint) *releaseInfo {
	rel := &shipper.Release{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "Release",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      incumbentName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion: shipper.SchemeGroupVersion.String(),
					Kind:       "Application",
					Name:       app.GetName(),
					UID:        app.GetUID(),
				},
			},
			Labels: map[string]string{
				shipper.ReleaseLabel: incumbentName,
				shipper.AppLabel:     app.GetName(),
			},
			Annotations: map[string]string{
				shipper.ReleaseGenerationAnnotation: "0",
			},
		},
		Status: shipper.ReleaseStatus{
			AchievedStep: &shipper.AchievedStep{
				Step: 2,
				Name: "full on",
			},
			Conditions: []shipper.ReleaseCondition{
				{Type: shipper.ReleaseConditionTypeComplete, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.ReleaseSpec{
			TargetStep: 2,
			Environment: shipper.ReleaseEnvironment{
				Strategy: &vanguard,
			},
		},
	}

	installationTarget := &shipper.InstallationTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "InstallationTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      incumbentName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.InstallationTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.InstallationTargetSpec{
			Clusters: []string{clusterName},
		},
	}

	capacityTarget := &shipper.CapacityTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "CapacityTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      incumbentName,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.CapacityTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.CapacityTargetSpec{
			Clusters: []shipper.ClusterCapacityTarget{
				{
					Name:              clusterName,
					Percent:           100,
					TotalReplicaCount: int32(totalReplicaCount),
				},
			},
		},
	}

	trafficTarget := &shipper.TrafficTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "TrafficTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      incumbentName,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.TrafficTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.TrafficTargetSpec{
			Clusters: []shipper.ClusterTrafficTarget{
				{
					Name:   clusterName,
					Weight: 100,
				},
			},
		},
	}

	return &releaseInfo{
		release:            rel,
		installationTarget: installationTarget,
		capacityTarget:     capacityTarget,
		trafficTarget:      trafficTarget,
	}
}

// buildContender returns a releaseInfo with a contender release and
// associated objects.
func buildContender(totalReplicaCount uint) *releaseInfo {
	rel := &shipper.Release{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "Release",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      contenderName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					APIVersion: shipper.SchemeGroupVersion.String(),
					Kind:       "Application",
					Name:       app.GetName(),
					UID:        app.GetUID(),
				},
			},
			Labels: map[string]string{
				shipper.ReleaseLabel: contenderName,
				shipper.AppLabel:     app.GetName(),
			},
			Annotations: map[string]string{
				shipper.ReleaseGenerationAnnotation: "1",
			},
		},
		Status: shipper.ReleaseStatus{
			Conditions: []shipper.ReleaseCondition{
				{Type: shipper.ReleaseConditionTypeScheduled, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.ReleaseSpec{
			TargetStep: 0,
			Environment: shipper.ReleaseEnvironment{
				Strategy: &vanguard,
			},
		},
	}

	installationTarget := &shipper.InstallationTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "InstallationTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      contenderName,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.InstallationTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.InstallationTargetSpec{
			Clusters: []string{clusterName},
		},
	}

	capacityTarget := &shipper.CapacityTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "CapacityTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      contenderName,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.CapacityTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.CapacityTargetSpec{
			Clusters: []shipper.ClusterCapacityTarget{
				{
					Name:              clusterName,
					Percent:           0,
					TotalReplicaCount: int32(totalReplicaCount),
				},
			},
		},
	}

	trafficTarget := &shipper.TrafficTarget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: shipper.SchemeGroupVersion.String(),
			Kind:       "TrafficTarget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      contenderName,
			OwnerReferences: []metav1.OwnerReference{metav1.OwnerReference{
				APIVersion: shipper.SchemeGroupVersion.String(),
				Name:       rel.Name,
				Kind:       "Release",
				UID:        rel.UID,
			}},
		},
		Status: shipper.TrafficTargetStatus{
			Conditions: []shipper.TargetCondition{
				{Type: shipper.TargetConditionTypeReady, Status: corev1.ConditionTrue},
			},
		},
		Spec: shipper.TrafficTargetSpec{
			Clusters: []shipper.ClusterTrafficTarget{
				{
					Name:   clusterName,
					Weight: 0,
				},
			},
		},
	}

	return &releaseInfo{
		release:            rel,
		installationTarget: installationTarget,
		capacityTarget:     capacityTarget,
		trafficTarget:      trafficTarget,
	}
}

// ensureCapacityPatch executes the strategy and returns the content of the
// patch it produced, if any. Returns an error if the number of patches is
// wrong, the patch doesn't implement the CapacityTargetOutdatedResult
// interface or the name is different than the given expectedName.
func ensureCapacityPatch(e *Executor, expectedName string, role role) (*shipper.CapacityTargetSpec, error) {
	if patches, _, err := e.Execute(); err != nil {
		return nil, err
	} else {

		aType := shipper.StrategyConditionContenderAchievedCapacity
		releaseRole := "contender"

		if role == Incumbent {
			aType = shipper.StrategyConditionIncumbentAchievedCapacity
			releaseRole = "incumbent"
		}

		strategyConditions := firstConditionsFromPatches(patches)

		if strategyConditions == nil {
			return nil, fmt.Errorf("could not find a ReleaseUpdateResult patch")
		}

		if s, ok := strategyConditions.GetStatus(aType); ok && s != corev1.ConditionFalse {
			return nil, fmt.Errorf("%s shouldn't have achieved traffic", releaseRole)
		}

		if len(patches) < 2 || len(patches) > 3 {
			return nil, fmt.Errorf("expected two or three patches, got %d patch(es) instead", len(patches))
		} else {
			if p, ok := patches[0].(*CapacityTargetOutdatedResult); !ok {
				return nil, fmt.Errorf("expected a CapacityTargetOutdatedResult, got something else")
			} else {
				if p.Name != expectedName {
					return nil, fmt.Errorf("expected %q, got %q", expectedName, p.Name)
				} else {
					return p.NewSpec, nil
				}
			}
		}
	}
}

// ensureTrafficPatch executes the strategy and returns the content of the
// patch it produced, if any. Returns an error if the number of patches is
// wrong, the patch doesn't implement the TrafficTargetOutdatedResult
// interface or the name is different than the given expectedName.
func ensureTrafficPatch(e *Executor, expectedName string, role role) (*shipper.TrafficTargetSpec, error) {
	if patches, _, err := e.Execute(); err != nil {
		return nil, err
	} else {

		aType := shipper.StrategyConditionContenderAchievedTraffic
		releaseRole := "contender"

		if role == Incumbent {
			aType = shipper.StrategyConditionIncumbentAchievedTraffic
			releaseRole = "incumbent"
		}

		strategyConditions := firstConditionsFromPatches(patches)

		if s, ok := strategyConditions.GetStatus(aType); ok && s != corev1.ConditionFalse {
			return nil, fmt.Errorf("%s shouldn't have achieved traffic", releaseRole)
		}

		if len(patches) != 2 {
			return nil, fmt.Errorf("expected two patches, got %d patch(es) instead", len(patches))
		} else {
			if p, ok := patches[0].(*TrafficTargetOutdatedResult); !ok {
				return nil, fmt.Errorf("expected a TrafficTargetOutdatedResult, got something else")
			} else {
				if p.Name != expectedName {
					return nil, fmt.Errorf("expected %q, got %q", expectedName, p.Name)
				} else {
					return p.NewSpec, nil
				}
			}
		}
	}
}

// ensureReleasePatch executes the strategy and returns the content of the
// patch it produced, if any. Returns an error if the number of patches is
// wrong, the patch doesn't implement the ReleaseUpdateResult interface or
// the name is different than the given expectedName.
func ensureReleasePatch(e *Executor, expectedName string) (*shipper.ReleaseStatus, error) {
	if patches, _, err := e.Execute(); err != nil {
		return nil, err
	} else {

		strategyConditions := firstConditionsFromPatches(patches)

		if !strategyConditions.AllTrue(e.contender.release.Spec.TargetStep) {
			return nil, fmt.Errorf("all conditions should have been met")
		}

		if len(patches) != 1 {
			return nil, fmt.Errorf("expected one patch, got %d patches instead", len(patches))
		} else {
			if p, ok := patches[0].(*ReleaseUpdateResult); !ok {
				return nil, fmt.Errorf("expected a ReleaseUpdateResult, got something else")
			} else {
				if p.Name != expectedName {
					return nil, fmt.Errorf("expected %q, got %q", expectedName, p.Name)
				} else {
					return p.NewStatus, nil
				}
			}
		}
	}
}

// ensureFinalReleasePatches executes the strategy and returns an error if
// the number of patches is wrong, the patches' phases are wrong for either
// incumbent or contender.
func ensureFinalReleasePatches(e *Executor) error {
	if patches, _, err := e.Execute(); err != nil {
		return err
	} else {

		strategyConditions := firstConditionsFromPatches(patches)

		if !strategyConditions.AllTrue(int32(e.contender.release.Spec.TargetStep)) {
			return fmt.Errorf("all conditions should be true")
		}

		if len(patches) != 1 {
			return fmt.Errorf("expected one patches, got %d patches instead", len(patches))
		} else {
			for _, patch := range patches {
				if p, ok := patch.(*ReleaseUpdateResult); !ok {
					return fmt.Errorf("expected a ReleaseUpdateResult, got something else")
				} else {
					if p.Name == contenderName {
						if p.NewStatus.AchievedStep == nil || p.NewStatus.AchievedStep.Step != 2 {
							return fmt.Errorf(
								"expected contender achievedSteps %d, got %v",
								2, p.NewStatus.AchievedStep)
						}
					}
				}
			}
		}
	}
	return nil
}

func firstConditionsFromPatches(patches []ExecutorResult) conditions.StrategyConditionsMap {
	var strategyConditions conditions.StrategyConditionsMap
	for p := range patches {
		if patch, ok := patches[p].(*ReleaseUpdateResult); ok {
			if patch.NewStatus.Strategy != nil {
				strategyConditions = conditions.NewStrategyConditions(patch.NewStatus.Strategy.Conditions...)
			}
		}
	}
	return strategyConditions
}
