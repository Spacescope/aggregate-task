package utils

// func FindIntersections(r [][]uint64) []uint64 {
// 	if len(r) <= 1 {
// 		return []uint64{}
// 	}

// 	result := r[0]
// 	for _, set := range r {
// 		result = SliceIntersectionInt(result, set)
// 	}

// 	return result
// }

// func SliceIntersectionInt(A, B []uint64) []uint64 {
// 	if len(A) < 1 || len(B) < 1 {
// 		return []uint64{}
// 	}
// 	result := make([]uint64, 0)

// 	// deduplication
// 	flagMap := make(map[uint64]bool, 0)
// 	for _, a := range A {
// 		if _, ok := flagMap[a]; ok {
// 			continue
// 		}
// 		flagMap[a] = true
// 		for _, b := range B {
// 			if b == a {
// 				result = append(result, a)
// 				break
// 			}
// 		}
// 	}
// 	return result
// }

// func SliceDifferenceInt(A, B []uint64) []uint64 {
// 	if len(A) < 1 || len(B) < 1 {
// 		return A
// 	}
// 	result := make([]uint64, 0)

// 	// deduplication
// 	flagMap := make(map[uint64]bool, 0)
// 	for _, a := range A {
// 		if _, ok := flagMap[a]; ok {
// 			continue
// 		}
// 		flagMap[a] = true
// 		flag := true
// 		for _, b := range B {
// 			if b == a {
// 				flag = false
// 				break
// 			}
// 		}
// 		if flag {
// 			result = append(result, a)
// 		}
// 	}
// 	return result
// }

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
