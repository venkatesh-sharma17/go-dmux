package core

//Hasher interface that can be implemented to define data interface and compute
// and return the righ hash int for it
type Hasher interface {
	//ComputeHash method has to be implemented by Client to define how to distribute
	// when using HashDistributor
	ComputeHash(data interface{}) int
}

//GetHashDistribution returns hashDistributor implementation of Distributor
//interface to provide Consistent Hash based routing in Dmux from Source to Sink.
//This needs Client to implement Hasher and pass Hasher in this method arg
func GetHashDistribution(h Hasher) Distributor {
	hd := hashDistributor{h}
	return &hd
}

//DistributorType based on this distribution is determined
type DistributorType string

const (
	//HashDistributor will distribute based on the key hash
	HashDistributor DistributorType = "Hash"
	//RoundRobinDistributor will distribute on round robin fashion
	RoundRobinDistributor DistributorType = "RoundRobin"
)

//GetDistribution returns correct Distributor based on distributorType
func GetDistribution(distributorType DistributorType, h Hasher) Distributor {
	switch distributorType {
	case RoundRobinDistributor:
		return GetRoundRobinDistribution()
	default:
		return GetHashDistribution(h)
	}
}

//GetRoundRobinDistribution return roundRobinDistributor
func GetRoundRobinDistribution() Distributor {
	hd := roundRobinDistributor{0}
	return &hd
}

type hashDistributor struct {
	hasher Hasher
}

type roundRobinDistributor struct {
	count int
}

func (h *roundRobinDistributor) Distribute(data interface{}, size int) int {
	count := h.count % size
	h.count = count + 1
	return count
}

func (h *hashDistributor) Distribute(data interface{}, size int) int {
	hash := h.hasher.ComputeHash(data)
	hash = abs(hash)
	bucket := mod(hash, size)
	// log.Printf("hash_bucket_choosen %d", bucket)
	return bucket
}

func abs(hash int) int {
	if hash < 0 {
		return -hash
	}
	return hash
}

func mod(hash, size int) int {
	val := hash % size // TODO change this to bit shiffing when size ^2
	return val
}
