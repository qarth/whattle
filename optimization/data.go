package optimization

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"

	log "github.com/cihub/seelog"
)

type (
	Data struct {
		Grid    `json:"grid"`
		EbvCols int         `json:"ebv_column"`
		Ebv     [][]float64 `json:"-"`
	}
)

func (block *Data) initializeFromGzip(infile string) error {

	f, e := os.Open(infile)

	if e != nil {
		log.Errorf("Error: failed initializing data from input file (path err) %v: %v", infile, e)
		return e
	}
	defer f.Close()

	r, e := gzip.NewReader(f)
	if e != nil {
		log.Errorf("Error: failed initializing data from input file (gzip err) %v: %v", infile, e)
		return e
	}
	defer r.Close()

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	//-------------------------------

	cnt := block.Grid.gridCount()
	block.Ebv = [][]float64{}
	idx := 0
	face := make([]float64, cnt)

	// Read every line
	for s.Scan() {
		v, e := strconv.ParseFloat(s.Text(), 64)
		if e != nil {
			return e
		}
		face[idx] = v

		// one layer has been read,begin next layer
		if idx++; idx >= cnt {
			layer := make([]float64, cnt)
			copy(layer, face)
			block.Ebv = append(block.Ebv, layer)
			idx = 0
		}
	}

	//-------------------------------

	if idx != 0 {
		e = fmt.Errorf("Error: failed initializing data from input file %v: Invalid input data file (idx not equal to 0)", infile)
	} else if len(block.Ebv) == 0 {
		e = fmt.Errorf("ERROR: no data")
	} else if len(block.Ebv[0]) == 0 {
		e = fmt.Errorf("ERROR: no values")
	} else if len(block.Ebv[0]) != cnt {
		e = fmt.Errorf("ERROR: wrong number of values")
	} else {
		e = nil
	}

	if e != nil {
		log.Error(e)
	}

	return e
}
