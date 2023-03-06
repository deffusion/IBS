package main

//func TestInfoSorter(t *testing.T) {
//	s := []int64{1, 3, 9, 5, 3, 6, 8, 1, 0, 1, 3, 6, 100, 34, 36, 67, 87, 33, 10, 4, 5, 8, 4, 9, 0, 7, 5, 4, 3, 2, 3, 5, 6, 7, 3, 1}
//	var infos []*information.Packet
//	for i := 0; i < len(s); i++ {
//		p := &information.Packet{}
//		p.Timestamp = s[i]
//		infos = append(infos, p)
//	}
//	sorter := NewInfoSorter()
//	for _, info := range infos {
//		fmt.Println("append", info.Timestamp)
//		sorter.Append(info)
//	}
//	sorter.Print()
//	l := sorter.length
//	for i := 0; i < l-1; i++ {
//		info, _ := sorter.Take()
//		fmt.Println("take", info.Timestamp)
//	}
//}
