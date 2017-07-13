package ecs

import (
	"sort"
	"sync"
	"sync/atomic"
)

var (
	counterLock sync.Mutex
	idInc       uint64
)

type Entity interface {
	ID() uint64
}

type entity uint64

type IdentifierSlice []Entity

func NewEntity() entity {
	return entity(atomic.AddUint64(&idInc, 1))
}

func NewEntitys(amount int) []Entity {
	entities := make([]Entity, amount)

	counterLock.Lock()
	for i := 0; i < amount; i++ {
		idInc++
		entities[i] = entity(idInc)
	}
	counterLock.Unlock()

	return entities
}

func (e entity) ID() uint64 {
	return uint64(e)
}

func (is IdentifierSlice) Len() int {
	return len(is)
}

func (is IdentifierSlice) Less(i, j int) bool {
	return is[i].ID() < is[j].ID()
}

func (is IdentifierSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

type Prioritizer interface {
	Priority() int
}

type System interface {
	Prioritizer
	Update(int64) error
	Remove(uint64)
}

type systems []System

func (s systems) Len() int {
	return len(s)
}

func (s systems) Less(i, j int) bool {
	var prio1, prio2 int
	prio1 = s[i].Priority()
	prio2 = s[j].Priority()
	return prio1 > prio2
}

func (s systems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type HandleErrorFn func(error)

type World struct {
	hefn    HandleErrorFn
	systems systems
}

func New(hefn HandleErrorFn) *World {
	return &World{
		hefn,
		make(systems, 0),
	}
}

func (w *World) Add(s ...System) {
	for _, sys := range s {
		w.systems = append(w.systems, sys)
	}
	sort.Sort(w.systems)
}

func (w *World) Systems() []System {
	return w.systems
}

func (w *World) Update(dt int64) {
	var err error
	for _, system := range w.Systems() {
		err = system.Update(dt)
		if err != nil {
			w.hefn(err)
		}
	}
}

func (w *World) RemoveEntity(e uint64) {
	for _, sys := range w.systems {
		sys.Remove(e)
	}
}
