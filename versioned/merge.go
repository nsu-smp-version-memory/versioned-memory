package versioned

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

func (s Set) Merge(vm *core.VersionManager, b Set) Set {
	base := core.CommonAncestor(s.version, b.version)
	version_parent := vm.NewChild(base)

	items := mergeKeyItems(s.root, b.root)
	root := buildBalanced(items)

	return Set{
		root:    root,
		version: version_parent,
	}
}

type keyItem struct {
	value int
	adds  TagSet
	dels  TagSet
}

func mergeKeyItems(a, b *node) []keyItem {
	var sa, sb []keyItem
	inOrderItems(a, &sa)
	inOrderItems(b, &sb)

	out := make([]keyItem, 0, len(sa)+len(sb))
	i, j := 0, 0

	for i < len(sa) && j < len(sb) {
		if sa[i].value < sb[j].value {
			out = append(out, sa[i])
			i++
		} else if sa[i].value > sb[j].value {
			out = append(out, sb[j])
			j++
		} else {
			out = append(out, keyItem{
				value: sa[i].value,
				adds:  unionTags(sa[i].adds, sb[j].adds),
				dels:  unionTags(sa[i].dels, sb[j].dels),
			})
			i++
			j++
		}
	}

	for i < len(sa) {
		out = append(out, sa[i])
		i++
	}
	for j < len(sb) {
		out = append(out, sb[j])
		j++
	}

	return out
}

func inOrderItems(n *node, out *[]keyItem) {
	if n == nil {
		return
	}
	inOrderItems(n.left, out)
	*out = append(*out, keyItem{
		value: n.value,
		adds:  n.adds,
		dels:  n.dels,
	})
	inOrderItems(n.right, out)
}

func unionTags(a, b TagSet) TagSet {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}
	out := make(TagSet, len(a)+len(b))
	for k := range a {
		out[k] = struct{}{}
	}
	for k := range b {
		out[k] = struct{}{}
	}
	return out
}

func buildBalanced(items []keyItem) *node {
	if len(items) == 0 {
		return nil
	}
	m := len(items) / 2
	return &node{
		value: items[m].value,
		adds:  items[m].adds,
		dels:  items[m].dels,
		left:  buildBalanced(items[:m]),
		right: buildBalanced(items[m+1:]),
	}
}
