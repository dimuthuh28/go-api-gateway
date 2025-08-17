package loadbalancer

type RoundRobin struct {
	backends []string
	counter  int
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{backends: backends}
}

func (r *RoundRobin) Next() string {
	backend := r.backends[r.counter%len(r.backends)]
	r.counter++
	return backend
}
