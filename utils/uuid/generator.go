package uuid

import (
	"sync"

	"math/rand"
)

// NewGenerator returns a new generator
func NewGenerator() *Generator {
	var g Generator
	// Set random from a newly generated seed
	g.rnd = rand.New(rand.NewSource(generateSeed()))
	return &g
}

// Generator is a thread-safe UUID generator
type Generator struct {
	mu  sync.Mutex
	rnd *rand.Rand
}

// New returns a new UUID
func (g *Generator) New() *UUID {
	u := g.Make()
	return &u
}

// Make initializes a UUID
func (g *Generator) Make() (u UUID) {
	// Acquire lock
	g.mu.Lock()
	// Create a new UUID from a random Int64 value
	// Note: This value is actually int63, but that's OK because newUUID truncates this to int48.
	u = makeUUID(g.rnd.Int63())
	// Release lock
	g.mu.Unlock()
	return
}
