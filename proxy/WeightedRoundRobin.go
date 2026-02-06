package proxy


// not a very intuitive idea. But i think its very interesting and more optimized and less complex 

// and adminapi make sure its correct, should i do the add with weight or not (initialize it always with 1)
// also in main i should see the new way of adding (genre la commande commentee dial curl POST the new version since now we add hta weight)
// IF I DO THIS? I SHOULD TEST AND BE DONE WITH THIS AND PROBABLY MOVE TO ANOTHER ENHACEMENT
type WeightedRoundRobin struct {
	Pool *ServerPool
}

func NewWeightedRoundRobin(pool *ServerPool) *WeightedRoundRobin {
	return &WeightedRoundRobin{Pool: pool}
}

func (wrr *WeightedRoundRobin) GetNextValidPeer() *Backend {
	wrr.Pool.mux.Lock()
	defer wrr.Pool.mux.Unlock()

	var selected *Backend
	totalWeight := 0

	// Step 1: Increase currentWeight for each alive backend and calculate total
	for _, backend := range wrr.Pool.Backends {
		if !backend.IsAlive() {
			continue
		}

		// Increase current weight by static weight
		backend.UpdateCurrentWeight(backend.Weight)
		totalWeight += backend.Weight

		// Select backend with highest current weight
		if selected == nil || backend.GetCurrentWeight() > selected.GetCurrentWeight() {
			selected = backend
		}
	}

	if selected == nil {
		return nil
	}

	selected.UpdateCurrentWeight(-totalWeight)

	return selected
}
