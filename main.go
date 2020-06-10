package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sausheong/petri"
)

var w int
var raceColors []int
var off = petri.White

// used to calculate the ratio of races for initial population
var raceRatio *string

// percentage of simulation grid that is populated
var coverage *float64

// minimum acceptable number of neighbours of same race
var minList *string

// maximum capped number of neighbours of the same race
var maxList *string

// min and max list in
var mins []int
var maxs []int

func init() {
	raceRatio = flag.String("ratio", "1:1", "Race ratio eg 1:1 or 1:2:2 etc. Numbers are positional")
	coverage = flag.Float64("coverage", 0.7, "percentage of simulation grid that is populated")
	minList = flag.String("min", "2:2", "minimum acceptable number of neighbours of same race eg 2:2 pr 3:2:4 etc. Numbers are positional")
	maxList = flag.String("max", "8:8", "max acceptable number of neighbours of same race eg 6:6 pr 6:7:8 etc. Numbers are positional")
	raceColors = []int{
		petri.Deeppink,
		petri.Lawngreen,
		petri.Deepskyblue,
		petri.Peachpuff,
		petri.Sandybrown,
		petri.Fuchsia,
	}
}

func main() {
	s := &Segregation{
		petri.Sim{},
	}
	petri.Run(s)
}

// Segregation is a simulation of racial segregation
type Segregation struct {
	petri.Sim
}

// Init creates the initial cell population
func (s *Segregation) Init() {
	w = *petri.Width
	rand.Seed(time.Now().UTC().UnixNano())
	s.Units = make([]petri.Cellular, w*w)
	pop := calc(*raceRatio)
	n := 0
	for i := 1; i <= w; i++ {
		for j := 1; j <= w; j++ {
			p, q := rand.Float64(), rand.Float64()
			if p < *coverage {
				c := raceColors[pop(q)]
				s.Units[n] = s.CreateCell(i, j, c, c)
			} else {
				s.Units[n] = s.CreateCell(i, j, off, off)
			}
			n++
		}
	}
	mins = split(*minList)
	maxs = split(*maxList)
	if len(mins) != len(maxs) || len(mins) != len(split(*raceRatio)) {
		fmt.Println("Lengths of minimum, maximum or population ratio must be the same, default is 2")
		os.Exit(0)
	}
}

// Process the simulation
func (s *Segregation) Process() {
	for cellNumber, cell := range s.Units {
		index := raceIndex(cell.RGB())
		if cell.RGB() == off {
			continue
		}
		// find all the cell's neighbours
		neighbours := petri.FindNeighboursIndex(cellNumber)
		// count of the neighbours that are the same as the cell
		sameCount := 0

		// for every neighbour
		for _, neighbour := range neighbours {
			// if the cell is empty, go to the next neighbour
			if s.Units[neighbour].RGB() == off {
				continue
			}
			// if the neighbour is the same, increment sameCount
			if s.Units[neighbour].RGB() == cell.RGB() {
				sameCount++
			}
		}

		// check min and max number of neighbours of same race
		if sameCount < mins[index] || sameCount > maxs[index] {
			// find an empty
			empty := s.findEmpty()
			e := s.findRandomEmpty(empty)
			// move the current cell to the empty cell
			s.Units[e].SetRGB(cell.RGB())
			s.Units[cellNumber].SetRGB(off)
		}
	}
}

// find the index of a random empty cell in the grid
func (s *Segregation) findRandomEmpty(empty []int) int {
	r := rand.Intn(len(empty))
	return empty[r]
}

// find all cells that are empty in the grid
func (s *Segregation) findEmpty() (empty []int) {
	for n, cell := range s.Units {
		if cell.RGB() == off {
			empty = append(empty, n)
		}
	}
	return
}

// split the list and convert it into integers
func split(l string) []int {
	str := strings.Split(l, ":")
	list := []int{}
	for _, i := range str {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err) // if the string are not integers
		}
		list = append(list, j)
	}
	return list
}

type population = func(a float64) int

// returns a function that calculates the population probabilities
// for populating the initial simulation
func calc(ratio string) population {
	str := strings.Split(ratio, ":")
	ratios := []int{}
	total := 0.0
	for _, i := range str {
		j, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println("Population ratio is not in correct format:", err)
			os.Exit(0)
		}
		ratios = append(ratios, j)
		total += float64(j)
	}
	var probabilities = []float64{}
	for _, i := range ratios {
		probabilities = append(probabilities, float64(i)/total)
	}

	return func(p float64) (ret int) {
		for i := range probabilities {
			sum := 0.0
			for a := 0; a < i+1; a++ {
				sum += probabilities[a]
			}
			if p < sum {
				ret = i
				break
			}
		}
		return
	}
}

// find the index of the race given the color
func raceIndex(race int) int {
	for i, r := range raceColors {
		if r == race {
			return i
		}
	}
	return -1
}
