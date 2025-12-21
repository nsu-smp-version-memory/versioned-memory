package set

import "sort"

func replayToMap(ops []operation) map[int]struct{} {
	if len(ops) == 0 {
		return map[int]struct{}{}
	}

	tmp := make([]operation, len(ops))
	copy(tmp, ops)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].id.Before(tmp[j].id) })

	present := make(map[int]struct{}, len(tmp))
	for _, o := range tmp {
		switch o.kind {
		case operationAdd:
			present[o.value] = struct{}{}
		case operationRemove:
			delete(present, o.value)
		}
	}
	return present
}

func mapKeysSorted(m map[int]struct{}) []int {
	values := make([]int, 0, len(m))
	for k := range m {
		values = append(values, k)
	}
	sort.Ints(values)
	return values
}
