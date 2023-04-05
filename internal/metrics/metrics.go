package metrics

var DefaultStorage *MemStorage

func init() {
	DefaultStorage = NewStorage()
}

type Metric interface {
	Update(interface{})
	Get() interface{}
}

type MemStorage struct {
	Gauges   map[string]Metric
	Counters map[string]Metric
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]Metric),
		Counters: make(map[string]Metric),
	}
}

type Gauge float64
type Counter int64

func (g *Gauge) Update(val interface{}) {
	*g = Gauge(val.(float64))
}

func (g *Gauge) Get() interface{} {
	return float64(*g)
}

func (c *Counter) Update(val interface{}) {
	*c += Counter(val.(int64))
}

func (c *Counter) Get() interface{} {
	return int64(*c)
}
