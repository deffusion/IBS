package num_set

type Set struct {
	s []uint64
}

func NewSet() *Set {
	return &Set{
		[]uint64{},
	}
}

func (t *Set) Insert(num uint64) bool {
	l := t.locate(num)
	//fmt.Println("t.s", t.s, num, l)
	if len(t.s) > 0 && t.s[l] == num {
		return false
	}
	t.s = append(t.s, num)
	for i := len(t.s) - 1; i > 0; i-- {
		if t.s[i-1] > t.s[i] {
			t.s[i], t.s[i-1] = t.s[i-1], t.s[i]
		}
	}
	return true
}

func (t *Set) locate(num uint64) int {
	left, right := 0, len(t.s)
	for left < right {
		mid := (left + right) / 2
		if num > t.s[mid] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	if left < len(t.s) {
		return left
	}
	return left - 1
}

// Around n number around
func (t *Set) Around(num uint64, n int) []uint64 {
	var set []uint64
	pos := t.locate(num)
	left := pos - 1
	right := pos + 1
	set = append(set, t.s[pos])
	n--
	for n > 0 {
		if n%2 == 0 && left >= 0 {
			set = append(set, t.s[left])
			left--
		} else if right < len(t.s) {
			set = append(set, t.s[right])
			right++
		} else {
			break
		}
		n--
	}
	return set
}
