package main

import(
	geoInt "github.com/tapglue/geohash" // Int
	"github.com/pierrre/geohash"
	"github.com/AlasdairF/Sort/Uint64"
	e "github.com/cinnabardk/allancorfix2/internal/errors"
	"github.com/dhconnelly/rtreego"
	"math"
	"github.com/AlasdairF/Conv"
)


func (_ Test) GeoHash() {
	lat := 56.162939
	long := 10.203921
	hash := geohash.EncodeAuto(lat, long)
	e.InfoLog.Println(hash)
}

// https://github.com/yinqiwen/ardb/blob/master/doc/spatial-index.md
// http://gis.stackexchange.com/questions/18330/would-it-be-possible-to-use-geohash-for-proximity-searches
func (_ Test) GeoHashInt() {
	lat := 56.162939
	long := 10.203921
	hash := geoInt.EncodeInt(lat, long, 52) //  0.5971 meters
	hash = geoInt.EncodeInt(lat, long, 18)  // 78 km
	e.InfoLog.Println(hash)

	neighbours := geoInt.EncodeNeighborsInt(hash, 18)
	neighbours = append(neighbours, hash)
	sortUint64.Asc(neighbours)
	e.InfoLog.Println(neighbours)
}

type Somewhere struct {
	location rtreego.Point
	name     string
	num      int
}

func (s *Somewhere) Bounds() *rtreego.Rect {
	// define the bounds of s to be a rectangle centered at s.location
	// with side lengths 2 * tol:
	return s.location.ToRect(0.01)
}
func (_ Test) RTree() {
	calcSetup()
	rt := rtreego.NewTree(2, 25, 50)

	points := []*Somewhere{
		{rtreego.Point{55.711311, 9.536354}, "Vejle", 1},
		{rtreego.Point{57.725004, 10.579186}, "Skagen", 1},
		{rtreego.Point{56.162939, 10.203921}, "Århus", 1},
		{rtreego.Point{56.315499, 10.317270}, "Hornslet", 2},
		{rtreego.Point{56.363358, 10.235109}, "Hønebjergvej 31", 3},
		{rtreego.Point{56.460584, 10.036539}, "Randers", 4},
		{rtreego.Point{56.037247, 9.929799}, "Skanderborg", 5},
		{rtreego.Point{56.176362, 9.554922}, "Silkeborg", 6},
		{rtreego.Point{55.728449, 9.112366}, "Billund", 7},
	}

	for i := range points {
		rt.Insert(points[i])
	}
	rt.Delete(points[0])

	// Get a slice of the k objects in rt closest to q:
	here := points[3]
	num := rt.Size()
	results := rt.NearestNeighbors(num, here.location) // num must be less or equal to the size of the RTree
	e.InfoLog.Println("\nHvad er nærmest", here.name, ":")
	for i := range results {
		there := results[i].(*Somewhere)
		distance := calcDist(here.location, there.location)
		rounded := RoundPlus(distance, 1)
		Int := int(rounded * 10)

		km, remainder := Int/10, Int%10
		str := conv.String(km) + "." + conv.String(remainder)

		e.InfoLog.Println(i, ": ", there.name, ". Afstand: ", str, "full: ", RoundPlus(distance, 4))
	}
}

var preCalcCos [2000]float64

func calcSetup() {
	fc := 0.0001 * 8500
	for i := 0; i < 2000; i++ {
		preCalcCos[i] = math.Cos(fc)
		fc += 0.0001 // 1 / 10.000
	}
}

// Pythagoras’ theorem on Equirectangular approximation
func calcDist(current, away []float64) float64 {

	const Rad = math.Pi / 180

	latR1 := current[0] * Rad
	longR1 := current[1] * Rad

	latR2 := away[0] * Rad
	longR2 := away[1] * Rad

	intCos := (uint32((latR1 + latR2) * 5000)) - 8500
	var cos float64
	if intCos < 2000 {
		cos = preCalcCos[intCos]
	} else {
		cos = math.Cos((latR1 + latR2) / 2)
	}
	x := (longR2 - longR1) * cos
	y := latR2 - latR1

	return math.Sqrt(x*x+y*y) * 6371.0
}

// Rounds to desired places
func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return roundFloat4(f*shift) / shift
}

// rounds to a number without decimals
func roundFloat4(a float64) float64 {
	if a < 0 {
		return math.Ceil(a - 0.5)
	}
	return math.Floor(a + 0.5)
}

func roundFloatFast(f float64, i int) float64 {
	tmp := int(f * 100)
	last := int(f*1000) - tmp*10
	if last >= 5 {
		tmp += 1
	}

	return float64(tmp) / 100
}

// return rounded version of x with prec precision.
func roundFloat(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	intermed += .5
	x = .5
	if frac < 0.0 {
		x = -.5
		intermed -= 1
	}
	if frac >= x {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / pow
}