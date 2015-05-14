package graphalgo

import "fmt"

type Stack []int

func (s *Stack) push(v int) {
	*s = append(*s, v)
}

func (s *Stack) pop() int {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

func (s *Stack) isEmpty() bool {
	if len(*s) == 0 {
		return true
	}
	return false
}

func (s *Stack) top() int {
	return (*s)[len(*s)-1]
}

func SCC_kosaraju(AdjacencyList [][]int, AdjacencyListTranspose [][]int, NodeCount int) []int {

	color := make([]int, NodeCount)
	order := make([]int, NodeCount)
	orderIndex := 0
	cheked := make([]int, NodeCount)

	dfsStack := Stack{}
	minUncolored := 0

	for minUncolored < len(AdjacencyList) {

		dfsStack.push(int(minUncolored))
		color[minUncolored] = 1

		for !dfsStack.isEmpty() {
			v := dfsStack.top()

			nextChild := nextChild(AdjacencyList, v, cheked, color)

			if nextChild == -1 {
				dfsStack.pop()
				order[orderIndex] = v
				orderIndex++
				color[v] = 2
			} else {
				dfsStack.push(nextChild)
				color[nextChild] = 1
			}

		}

		i := minUncolored
		flag := false
		for i < len(AdjacencyList) && !flag {
			if color[i] == 0 {
				flag = true
			} else {
				i++
			}
		}
		minUncolored = i

	}

	curComponentIndex := int(1)
	maxUncolored := len(order) - 1
	color2 := make([]int, NodeCount)
	cheked2 := make([]int, NodeCount)

	for maxUncolored > 0 {

		dfsStack.push(int(order[maxUncolored]))
		color2[order[maxUncolored]] = curComponentIndex

		for !dfsStack.isEmpty() {
			v := dfsStack.top()

			nextChild := nextChild(AdjacencyListTranspose, v, cheked2, color2)

			if nextChild == -1 {
				dfsStack.pop()
			} else {
				dfsStack.push(nextChild)
				color2[nextChild] = curComponentIndex
			}

		}

		curComponentIndex++

		i := maxUncolored
		flag := false
		for i >= 0 && !flag {
			if color2[order[i]] == 0 {
				flag = true
			} else {
				i--
			}
		}
		maxUncolored = i
	}

	return color2
}

func nextChild(AdjacencyList [][]int, v int, cheked []int, color []int) int {

	i := cheked[v]
	for {

		if v >= int(len(AdjacencyList)) {
			cheked[v] = i
			return -1
		}
		if i >= int(len(AdjacencyList[v])) {
			cheked[v] = i
			return -1
		}

		if color[AdjacencyList[v][i]] == 0 {
			cheked[v] = i
			return AdjacencyList[v][i]
		}
		i++
	}
}

func CountSCCFreqency(scc []int) []int {

	freq := make([]int, scc[MaxValueIndex(scc)]+1)

	for i, _ := range scc {
		freq[scc[i]]++
	}
	return freq
}

func MaxValueIndex(slice []int) int {

	max := slice[0]
	q := int(0)
	for i, _ := range slice {
		if slice[i] > max {
			max = slice[i]
			q = int(i)
		}
	}
	return q
}

func MaxSCC(scc []int, count int, maxIndex int) []int {
	mScc := make([]int, count)
	j := 0
	for i, v := range scc {
		if maxIndex == v {
			mScc[j] = i
			j++
		}
	}
	return mScc
}

func BTModel(AdjacencyList [][]int, AdjacencyListTranspose [][]int, NodeCount int) []int {

	// Slice, где ключ это id вершины в графе, а значение это номер компоненты сильной связности, которой эта вершина принадлежит
	scc := SCC_kosaraju(AdjacencyList, AdjacencyListTranspose, NodeCount)
	// Slice, где ключ это номер компоненты сильной связности, а значение это количество вершин в этой компоненте.
	fscc := CountSCCFreqency(scc)

	indexMaxSCC := MaxValueIndex(fscc) // Номер максимальной компоненты сильной связности.

	btm := make([]int, NodeCount)

	// 0 - все вершины
	// 1 - принадлежат максимальной компоненте сильной связности
	// 2 - принадлежат истокам(in)
	// 4 - принадлежат стокам (out)
	// 8 - принадлежат t in
	// 16 - принадлежат t out
	// 12, 18 - tubes
	for i := 0; i < NodeCount; i++ {
		btm[i] = 0
		if scc[i] == indexMaxSCC {
			btm[i] = 1
		}
	}

	for i := 0; i < NodeCount; i++ {
		if btm[i] == 1 {
			for _, v := range AdjacencyListTranspose[i] {
				if btm[v] != 1 {
					btm[v] = 2
				}
			}
		}
	}

	for i := 0; i < NodeCount; i++ {
		if btm[i] == 1 {
			for _, v := range AdjacencyList[i] {
				if btm[v] != 1 {
					btm[v] = 4
				}
			}
		}
	}

	for i := 0; i < NodeCount; i++ {
		if btm[i] == 2 {
			for _, v := range AdjacencyList[i] {
				if btm[v] == 0 || btm[v] == 4 {
					btm[v] = btm[v] + 8
				}
			}
		}
	}

	for i := 0; i < NodeCount; i++ {
		if btm[i] == 4 {
			for _, v := range AdjacencyListTranspose[i] {
				if btm[v] == 0 || btm[v] == 2 || btm[v] == 8 {
					btm[v] = btm[v] + 16
				}
			}
		}
	}

	return btm
}

func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(x), cap(x), x)
}
