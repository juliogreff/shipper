package rolloutblock

import (
	"github.com/golang/glog"
	"sort"
	"strings"

	shipper "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	shippererrors "github.com/bookingcom/shipper/pkg/errors"
)

// Return a boolean to decide if the application should override one of the
// existing rollout blocks, a string representing the existing rollout blocks
// that are not overridden, and an error if the application is trying to override
// a non-existing rolloutblock
func ShouldOverrideRolloutBlock(overrideRB string, RBs []*shipper.RolloutBlock) (bool, string, error) {
	overrideRBs := strings.Split(overrideRB, ",")
	if len(overrideRB) == 0 {
		overrideRBs = []string{}
	}

	nonOverriddenRBs, err := Difference(RBs, overrideRBs)
	sort.Slice(nonOverriddenRBs, func(i, j int) bool {
		return nonOverriddenRBs[i] < nonOverriddenRBs[j]
	})

	return len(nonOverriddenRBs) == 0, strings.Join(nonOverriddenRBs, ", "), err
}

// Set difference: existingRBs - overrideRBs
// Return all RolloutBlocks that are not overridden.
// Return an InvalidRolloutBlockOverrideError if trying to override a
// non-existing rolloutblock object
func Difference(existingRBs []*shipper.RolloutBlock, overrideRBs []string) ([]string, error) {
	var diff []string
	overrideRBdict := make(map[string]bool)
	existingRBdict := make(map[string]bool)

	for _, item := range overrideRBs {
		overrideRBdict[item] = true
	}

	// Collecting rolloutblocks that exist but are not overridden
	for _, item := range existingRBs {
		rbFullName := item.Namespace + "/" + item.Name
		existingRBdict[rbFullName] = true
		if _, ok := overrideRBdict[rbFullName]; !ok {
			diff = append(diff, rbFullName)
		}
	}

	// Non-existing blocks enlisted in this annotation should not be allowed.
	for _, item := range overrideRBs {
		if _, ok := existingRBdict[item]; !ok {
			glog.Infof("This rollout block %s does not exists", item)
			return diff, shippererrors.NewInvalidRolloutBlockOverrideError(item)
		}
	}

	return diff, nil
}
