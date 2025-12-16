package set

func (s Set) Has(key int) bool {
	cur := s.root
	for cur != nil {
		if key < cur.value {
			cur = cur.left
		} else if key > cur.value {
			cur = cur.right
		} else {
			for tag := range cur.adds {
				if _, removed := cur.dels[tag]; !removed {
					return true
				}
			}
			return false
		}
	}
	return false
}
