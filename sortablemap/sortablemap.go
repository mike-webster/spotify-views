package sortablemap

type SortableMapItem struct {
	Key   string
	Value int32
}
type SortableMap []SortableMapItem

func (p SortableMap) Len() int           { return len(p) }
func (p SortableMap) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p SortableMap) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SortableMap) Contains(key string) int {
	for i, ii := range p {
		if ii.Key == key {
			return i
		}
	}
	return -1
}
func (p SortableMap) ToMap() map[string]int32 {
	ret := map[string]int32{}
	for _, i := range p {
		ret[i.Key] = i.Value
	}
	return ret
}

func GetSortableMap32(m map[string]int32) SortableMap {
	ret := SortableMap{}
	for k, v := range m {
		ret = append(ret, SortableMapItem{Key: k, Value: v})
	}
	return ret
}

func GetSortableMap(m map[string]int) SortableMap {
	ret := SortableMap{}
	for k, v := range m {
		ret = append(ret, SortableMapItem{Key: k, Value: int32(v)})
	}
	return ret
}
