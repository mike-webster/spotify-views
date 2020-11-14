package spotify

// DEPRECATED - DO NOT USE THIS ANYMORE.
// ALTERNATIVE: sortablemap

// Pair is an outdated way to sort a map
type Pair struct {
	Key   string
	Value int32
}

// Pairs is a collection of Pair
type Pairs []Pair

func (p Pairs) Len() int           { return len(p) }
func (p Pairs) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p Pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Contains returns the index of the element in the collection, if it exists.
func (p Pairs) Contains(key string) int {
	for i, ii := range p {
		if ii.Key == key {
			return i
		}
	}
	return -1
}

// ToMap returns a map representation of the Pairs
func (p Pairs) ToMap() map[string]int32 {
	ret := map[string]int32{}
	for _, i := range p {
		ret[i.Key] = i.Value
	}
	return ret
}
