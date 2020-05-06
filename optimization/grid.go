package optimization

import "fmt"

type Grid struct {
	NumX int `json:"num_x"`
	NumY int `json:"num_y"`
	NumZ int `json:"num_z"`

	MinX float64 `json:"min_x"`
	MinY float64 `json:"min_y"`
	MinZ float64 `json:"min_z"`

	SizX float64 `json:"siz_x"`
	SizY float64 `json:"siz_y"`
	SizZ float64 `json:"siz_z"`

	gridcnt int
}

func (grid *Grid) adjust4gslib() {
	grid.MinX = grid.SizX / 2.0
	grid.MinY = grid.SizY / 2.0
	grid.MinZ = grid.SizZ / 2.0
}

/**
 * Params:
 *  k = The one dimensional Grid index.
 * Returns: The Grid index in the x direction
 */
func (grid *Grid) gridIx(k int) int {
	return (k % (grid.NumX * grid.NumY)) % grid.NumX
}

/**
 * Params:
 *  k = The one dimensional Grid index.
 * Returns: The Grid index in the y direction
 */
func (grid *Grid) gridIy(k int) int {
	return (k % (grid.NumX * grid.NumY)) / grid.NumX
}

/**
 * Params:
 *  k = The one dimensional Grid index.
 * Returns: The Grid index in the z direction
 */
func (grid *Grid) gridIz(k int) int {
	return k / (grid.NumX * grid.NumY)
}

/**
 * Params:
 *  ix = The Grid index in the x direction
 *  iy = The Grid index in the y direction
 *  iz = The Grid index in the z direction
 * Returns: The one dimensional Grid index.
 */
func (grid *Grid) gridIndex(ix, iy, iz int) int {
	return (ix + iy*grid.NumX + iz*grid.NumX*grid.NumY)
}

func (grid *Grid) gridIndex2(ids []int) int {
	return grid.gridIndex(ids[0], ids[1], ids[2])
}

/**
 * Params:
 *  k = The one dimensional Grid index.
 *  x = The test point x coordinate
 *  y = The test point y coordinate
 *  z = The test point z coordinate
 * Returns: True if x, y, z is within the block.
 */
func (grid *Grid) gridPointInCell(k int, x, y, z float64) bool {

	xn := x - (float64(grid.gridIx(k))*grid.SizX + grid.MinX)
	yn := y - (float64(grid.gridIy(k))*grid.SizY + grid.MinY)
	zn := z - (float64(grid.gridIz(k))*grid.SizZ + grid.MinZ)

	retval := ((0.0 <= xn) && (xn < grid.SizX))
	retval = retval && ((0.0 <= yn) && (yn < grid.SizY))
	retval = retval && ((0.0 <= zn) && (zn < grid.SizZ))

	return retval
}

// Write the Grid definition to standard output.
func (grid *Grid) String() string {
	retval := fmt.Sprintf("x =%7d %12.1f  %10.1f\n", grid.NumX, grid.MinX, grid.SizX)
	retval += fmt.Sprintf("y =%7d %12.1f  %10.1f\n", grid.NumY, grid.MinY, grid.SizY)
	retval += fmt.Sprintf("z =%7d %12.1f  %10.1f", grid.NumZ, grid.MinZ, grid.SizZ)
	return retval
}

// The number of blocks
func (grid *Grid) gridCount() int {
	if grid.gridcnt <= 0 {
		grid.gridcnt = grid.NumX * grid.NumY * grid.NumZ
	}
	return grid.gridcnt
}

// The Grid's bounding axis aligned bounding box
func (grid *Grid) aabb() [6]float64 {
	retval := [6]float64{
		grid.MinX,
		grid.MinY,
		grid.MinZ,
		grid.MinX + float64(grid.NumX)*grid.SizX,
		grid.MinY + float64(grid.NumY)*grid.SizY,
		grid.MinZ + float64(grid.NumZ)*grid.SizZ,
	}
	return retval
}

func (grid *Grid) blockAABB(k int) [6]float64 {

	centroid := grid.blockCentroid2(k)

	halfSizX := grid.SizX / 2.0
	halfSizY := grid.SizY / 2.0
	halfSizz := grid.SizZ / 2.0

	AABB := [6]float64{
		centroid[0] - halfSizX,
		centroid[1] - halfSizY,
		centroid[2] - halfSizz,
		centroid[0] + halfSizX,
		centroid[1] + halfSizY,
		centroid[2] + halfSizz,
	}

	return AABB
}

/**
 * Params:
 *  k = The one dimensional Grid index.
 * Returns: The centroid as [x, y, z]
 */
func (grid *Grid) blockCentroid(i, j, k int) [3]float64 {
	return [3]float64{
		float64(i)*grid.SizX + grid.MinX + grid.SizX/2.0,
		float64(j)*grid.SizY + grid.MinY + grid.SizY/2.0,
		float64(k)*grid.SizZ + grid.MinZ + grid.SizZ/2.0,
	}
}

func (grid *Grid) blockCentroid2(k int) [3]float64 {
	return grid.blockCentroid(grid.gridIx(k), grid.gridIy(k), grid.gridIz(k))
}
