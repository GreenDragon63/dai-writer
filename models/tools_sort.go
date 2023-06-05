package models

import (
	"os"
	"strconv"
	"strings"
)

type Numeric []os.FileInfo

func (nf Numeric) Len() int {
	return len(nf)
}

func (nf Numeric) Swap(i, j int) {
	nf[i], nf[j] = nf[j], nf[i]
}

func (nf Numeric) Less(i, j int) bool {
	pathA := nf[i].Name()
	pathB := nf[j].Name()
	a, err1 := strconv.ParseInt(pathA[0:strings.LastIndex(pathA, ".")], 10, 64)
	b, err2 := strconv.ParseInt(pathB[0:strings.LastIndex(pathB, ".")], 10, 64)
	if err1 != nil || err2 != nil {
		return pathA < pathB
	}
	return a < b
}
