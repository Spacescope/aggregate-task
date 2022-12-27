package utils

func Intersect(lists ...[]uint64) []uint64 {
	var inter []uint64
	mp := make(map[uint64]int)
	l := len(lists)

	if l == 0 {
		return make([]uint64, 0)
	}
	if l == 1 {
		for _, s := range lists[0] {
			if _, ok := mp[s]; !ok {
				mp[s] = 1
				inter = append(inter, s)
			}
		}
		return inter
	}

	for _, s := range lists[0] {
		if _, ok := mp[s]; !ok {
			mp[s] = 1
		}
	}

	for _, list := range lists[1 : l-1] {
		for _, s := range list {
			if _, ok := mp[s]; ok {
				mp[s]++
			}
		}
	}

	for _, s := range lists[l-1] {
		if _, ok := mp[s]; ok {
			if mp[s] == l-1 {
				inter = append(inter, s)
			}
		}
	}

	return inter
}

func Except(a []uint64, b []uint64) []uint64 {
	var inter []uint64
	mp := make(map[uint64]bool)

	for _, s := range a {
		mp[s] = true
	}

	for _, s := range b {
		if _, ok := mp[s]; ok {
			delete(mp, s)
		}
	}
	for key := range mp {
		inter = append(inter, key)
	}
	return inter
}
