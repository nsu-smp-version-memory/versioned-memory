package set

func (s Set) Merge(other Set) Set {
	a := s.ToSortedSlice()
	b := other.ToSortedSlice()

	merged := mergeUniqueSorted(a, b)
	root := buildBalanced(merged)

	return Set{
		root: root,
		size: len(merged),
	}
}

func mergeUniqueSorted(a, b []int) []int {
	out := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			out = append(out, a[i])
			i++
		} else if a[i] > b[j] {
			out = append(out, b[j])
			j++
		} else {
			out = append(out, a[i])
			i++
			j++
		}
	}

	for i < len(a) {
		out = append(out, a[i])
		i++
	}
	for j < len(b) {
		out = append(out, b[j])
		j++
	}

	return out
}

func buildBalanced(values []int) *node {
	if len(values) == 0 {
		return nil
	}
	mid := len(values) / 2
	return &node{
		value: values[mid],
		left:  buildBalanced(values[:mid]),
		right: buildBalanced(values[mid+1:]),
	}
}
