package scheduler

import (
	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"sort"
)

func getOfferScalar(offer *mesos.Offer, name string) float64 {
	resources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == name
	})

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}

func getOfferCpu(offer *mesos.Offer) float64 {
	return getOfferScalar(offer, "cpus")
}

func getOfferMem(offer *mesos.Offer) float64 {
	return getOfferScalar(offer, "mem")
}

func logOffers(offers []*mesos.Offer) {
	for _, offer := range offers {
		log.Infof("Received Offer <%v> with cpus=%v mem=%v", offer.Id.GetValue(), getOfferCpu(offer), getOfferMem(offer))
	}
}

// Sort Functions
type Slice struct {
	sort.Float64Slice
	idx []int
}


func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}


func NewSlice(n []float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
