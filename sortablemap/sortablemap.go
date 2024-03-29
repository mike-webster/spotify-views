package sortablemap

import "errors"

// Item represents an element in a SortableMap
type Item struct {
	Key   string
	Value int32
}

// Map is a sortable map
type Map []Item

func (m Map) Len() int           { return len(m) }
func (m Map) Less(i, j int) bool { return m[i].Value < m[j].Value }
func (m Map) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// Contains returns the index of the key if it exists in the map. If no match
// is found -1 is returned.
func (m Map) Contains(key string) int {
	for i, ii := range m {
		if ii.Key == key {
			return i
		}
	}
	return -1
}

// ToMap will return a basic map version of the SortableMap
func (m Map) ToMap() map[string]int32 {
	ret := map[string]int32{}
	for _, i := range m {
		ret[i.Key] = i.Value
	}
	return ret
}

// GetSortableMap takes a map and returns a map that can be
// sorted.
func GetSortableMap(m map[string]int) Map {
	ret := Map{}
	for k, v := range m {
		ret = append(ret, Item{Key: k, Value: int32(v)})
	}
	return ret
}

// Take will return a subset of the Map, with the top elements of
// the given number
func (m Map) Take(limit int) Map {
	ret := Map{}
	for i, ii := range m {
		if i >= limit {
			break
		}
		ret = append(ret, Item{Key: ii.Key, Value: ii.Value})
	}
	return ret
}

func (m Map) Value(k string) (int, error) {
	for _, i := range m {
		if i.Key == k {
			return int(i.Value), nil
		}
	}
	return 0, errors.New("404")
}

// Item represents an element in a SortableMap
type FloatItem struct {
	Key   string
	Value float32
}

// Map is a sortable map
type FloatMap []FloatItem

func (m FloatMap) Len() int           { return len(m) }
func (m FloatMap) Less(i, j int) bool { return m[i].Value < m[j].Value }
func (m FloatMap) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

// Contains returns the index of the key if it exists in the map. If no match
// is found -1 is returned.
func (m FloatMap) Contains(key string) int {
	for i, ii := range m {
		if ii.Key == key {
			return i
		}
	}
	return -1
}

// ToMap will return a basic map version of the SortableMap
func (m FloatMap) ToMap() map[string]float32 {
	ret := map[string]float32{}
	for _, i := range m {
		ret[i.Key] = i.Value
	}
	return ret
}

// GetSortableMap takes a map and returns a map that can be
// sorted.
func GetSortableFloatMap(m map[string]float32) FloatMap {
	ret := FloatMap{}
	for k, v := range m {
		ret = append(ret, FloatItem{Key: k, Value: float32(v)})
	}
	return ret
}

// Take will return a subset of the Map, with the top elements of
// the given number
func (m FloatMap) Take(limit int) FloatMap {
	ret := FloatMap{}
	for i, ii := range m {
		if i >= limit {
			break
		}
		ret = append(ret, FloatItem{Key: ii.Key, Value: ii.Value})
	}
	return ret
}
