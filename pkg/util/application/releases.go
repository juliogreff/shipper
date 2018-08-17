package application

import (
	releaseutil "github.com/bookingcom/shipper/pkg/util/release"
	shipperV1 "github.com/bookingcom/shipper/pkg/apis/shipper/v1"
	"github.com/bookingcom/shipper/pkg/errors"
	"fmt"
	"sort"
)

func GetContender(appName string, rels []*shipperV1.Release) (*shipperV1.Release, error) {
	if len(rels) == 0 {
		return nil, errors.NewContenderNotFoundError(appName)
	}
	return rels[0], nil
}

func GetIncumbent(appName string, rels []*shipperV1.Release) (*shipperV1.Release, error) {
	for _, r := range rels {
		if releaseutil.ReleaseComplete(r) {
			return r, nil
		}
	}
	return nil, errors.NewIncumbentNotFoundError(appName)
}

func SortReleases(rels []*shipperV1.Release) ([]*shipperV1.Release, error) {
	type releaseAndGeneration struct {
		release    *shipperV1.Release
		generation int
	}

	filteredRels := make([]releaseAndGeneration, 0, len(rels))
	for _, rel := range rels {
		if rel.DeletionTimestamp != nil {
			continue
		}
		g, err := releaseutil.GetGeneration(rel)
		if err != nil {
			return nil, fmt.Errorf(`incomplete Release "%s/%s": %s`, rel.Namespace, rel.Name, err)
		}
		filteredRels = append(filteredRels, releaseAndGeneration{rel, g})
	}

	sort.Slice(filteredRels, func(i, j int) bool {
		return filteredRels[i].generation > filteredRels[j].generation
	})

	relsToReturn := make([]*shipperV1.Release, 0, len(filteredRels))
	for _, e := range filteredRels {
		relsToReturn = append(relsToReturn, e.release)
	}

	return relsToReturn, nil
}

func ReleasesToHistory(releases []*shipperV1.Release) []string {
	sort.Slice(releases, func(i, j int) bool {
		gi, _ := releaseutil.GetGeneration(releases[i])
		gj, _ := releaseutil.GetGeneration(releases[j])
		return gi < gj
	})
	names := make([]string, 0, len(releases))
	for _, rel := range releases {
		names = append(names, rel.GetName())
	}
	return names

}