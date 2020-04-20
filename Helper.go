package main

func IsBetween(id1 uint64, id2 uint64, id3 uint64, includeIs3 bool) bool {
	if includeIs3 && id2 == id3 {
		return true
	} else if id2 == id1 || id2 == id3 {
		return false
	}
	//diff21 := id2 - id1
	//diff32 := id3 - id2
	//diff13 := id1 - id3
	//return (diff21 + diff32) < diff13
	if id1 < id3 {
		return id1 < id2 && id2 < id3
	} else {
		return !(id3 < id2 && id2 < id1)
	}
}

func AddPow(id uint64, next int) uint64 {
	toAdd := uint64(0)
	for i := 0; i < next - 1; i++ {
		toAdd *= uint64(2)
	}
	return id + toAdd
}