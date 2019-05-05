module utils

/* Each block is identified in two ways: Using three indices
ix, iy, and iz which are indices in the x, y, and z directions,
respectively, or by using a single index idx which counts first
along x, then y, and then z */

type (
	Block struct {
		ix  uint
		iy  uint
		iz  uint
		ebv float32
	}
)

// type (
// 	Data struct {
// 		Grid    `json:"grid"`
// 		EbvCols int         `json:"ebv_column"`
// 		Ebv     [][]float64 `json:"-"`
// 	}
// )

var idx []int64

func idx2mca(idx []int64) error {
	for i, v := range idx {
		return
	}
}
