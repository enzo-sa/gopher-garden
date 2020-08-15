package engine

import (
	"gioui.org/f32"
	"github.com/enzo-sa/gopher-garden/quickrand"
)

// all game logic is handled in engine

const (
	Holes      = 5
	Snakes     = 3
	Carrots    = 2
	LawnLength = 7
	LawnArea   = LawnLength * LawnLength
)

// garden stores main info for game to be played
type Garden struct {
	Lawn  *[LawnArea]Grass
	ind   indeces
	Score int
	Dead  bool
	// prev_hole stores the previous hole the gopher was in before a swap
	prev_hole int // -1 if gopher has not just made a hole swap
}

// stores indeces relative to Lawn[] of each important
// game piece
type indeces struct {
	player  int
	holes   [Holes]int
	snakes  [Snakes]int
	carrots [Carrots]int
}

// individual elements Lawn is composed of
type Grass struct {
	Off    f32.Point
	Snake  entity
	Player entity
	Carrot bool
	Hole   bool
}

type entity struct {
	Has   bool
	Direc int // 0-up, 1-down, 2-right, 3-left
}

// initialize new Garden for new game
func NewGame() *Garden {
	var ga Garden
	var Lawn [LawnArea]Grass
	// offsets are initialized from a ui call to the engine
	// get distinct random seed values and initialize garden pieces from them
	seeds := quickrand.RandInts(LawnArea-1, Holes+Carrots+Snakes+1, true)
	for i, seed := range *seeds {
		if i == 0 {
			ga.ind.player = seed
			Lawn[seed].Player.Has = true
		} else if i > 0 && i <= Holes {
			ga.ind.holes[i-1] = seed
			Lawn[seed].Hole = true
		} else if i > Holes && i <= Holes+Snakes {
			ga.ind.snakes[i-Holes-1] = seed
			Lawn[seed].Snake.Has = true
		} else {
			ga.ind.carrots[i-Holes-Snakes-1] = seed
			Lawn[seed].Carrot = true
		}
	}
	ga.Lawn = &Lawn
	return &ga
}

// key handling for gopher movement
// returns wether update was caused by the key
func (ga *Garden) HandleKey(key string) bool {
	// e is 1000 so when testing inc it fails rather than adding an extra conditional
	var inc_map = map[string]int{"W": -LawnLength, "A": -1, "S": LawnLength, "D": 1, "E": 1000,
		"↑": -LawnLength, "←": -1, "↓": LawnLength, "→": 1, "Space": 1000, " ": 1000}
	var direc_map = map[string]int{"W": 0, "A": 3, "S": 1, "D": 2, "E": 0,
		"↑": 0, "←": 3, "↓": 1, "→": 2, "Space": 0, " ": 0}
	if inc, ok := inc_map[key]; ok {
		in_val_arg := ga.ind.holes[:]
		if (key == "E" || key == "Space" || key == " ") && quickrand.InVals(ga.ind.player, &in_val_arg) && ga.prev_hole == -1 {
			ga.prev_hole = ga.ind.player
			ga.Lawn[ga.ind.player].Player.Has = false
			// make sure gopher doesn't end up at same hole
			temp := *quickrand.RandInts(len(ga.ind.holes)-1, 2, true)
			if ga.ind.player == ga.ind.holes[temp[0]] {
				ga.ind.player = ga.ind.holes[temp[1]]
			} else {
				ga.ind.player = ga.ind.holes[temp[0]]
			}
			ga.Lawn[ga.ind.player].Player.Has = true
			ga.Lawn[ga.ind.player].Player.Direc = 0
			return true
		} else {
			// handle basic gopher movement
			if ga.ind.player+inc >= 0 && ga.ind.player+inc < LawnArea &&
				(direc_map[key] <= 1 || ga.ind.player/LawnLength == (ga.ind.player+inc)/LawnLength) {
				ga.Lawn[ga.ind.player].Player.Has = false
				ga.ind.player += inc
				ga.Lawn[ga.ind.player].Player.Has = true
				ga.Lawn[ga.ind.player].Player.Direc = direc_map[key]
				ga.prev_hole = -1
				return true
			}
		}
	}
	return false
}

// returns slice of all entity and carrot vacant Grass inds
func (ga *Garden) vacantGrass() *[]int {
	var vacant []int
	for i := 0; i < LawnArea; i++ {
		if !ga.Lawn[i].Player.Has && !ga.Lawn[i].Snake.Has && !ga.Lawn[i].Carrot {
			vacant = append(vacant, i)
		}
	}
	return &vacant
}

// checks if there are no carrots left. if so, set their
// position to Carrots # of new spots (not occupied by entities) on the garden
func (ga *Garden) handleCarrots() {
	// make sure all carrots are gone
	for i := 0; i < Carrots; i++ {
		if ga.ind.carrots[i] != -1 {
			return
		}
	}
	// add new carrots in random vacant Grass squares
	vacant := ga.vacantGrass()
	if len(*vacant) < Carrots {
		panic("Not Enough Grass-Space for Carrots.")
	} else {
		// get distinct seeds to index vacant by for the carrot num
		seeds := quickrand.RandInts(len(*vacant)-1, Carrots, true)
		for i, seed := range *seeds {
			ga.ind.carrots[i] = (*vacant)[seed]
			ga.Lawn[(*vacant)[seed]].Carrot = true
		}
	}
}

// handle periodic snake movements
func (ga *Garden) MoveSnakes() {
	// FIX
	var inc_map = map[int]int{0: -LawnLength, 1: LawnLength, 2: 1, 3: -1}
	inBounds := func(old, new, direc int) bool {
		// check that snakes are in the garden bounds and remain on their row or column respectively
		if (direc > 1 && old/LawnLength != new/LawnLength) || new >= LawnArea || new < 0 || ga.Lawn[new].Snake.Has {
			return false
		}
		return true
	}
	// get new direc for snake that is valid, if none, panic (no directions will happen when the snakes
	// cluster and one gets trapped, and I am still working out how to handle thats, but once I figure it out it wont panic)
	getDirec := func(i, direc int) (int, bool) {
		var direcs []int
		// iterate over directions, if i == direc (direc that we want to change), don't check it
		for j := 0; j < 4; j++ {
			if j != direc {
				if inBounds(ga.ind.snakes[i], ga.ind.snakes[i]+inc_map[j], j) {
					direcs = append(direcs, j)
				}
			}
		}
		if len(direcs) == 0 {
			// check top row for the snake to come in through after tangling
			// this method of solving snake tanglings makes the most sense and shouldnt cause complications
			// as there should always be at least one open spot for a snake on any given row (2 carrots, 3 snakes, 1 player, 7 spots in row)
			for j := 0; j < LawnLength; j++ {
				if !ga.Lawn[j].Snake.Has && !ga.Lawn[j].Carrot && !ga.Lawn[j].Player.Has {
					ga.Lawn[ga.ind.snakes[i]].Snake.Has = false
					ga.ind.snakes[i] = j
					ga.Lawn[ga.ind.snakes[i]].Snake.Direc = 1
					ga.Lawn[ga.ind.snakes[i]].Snake.Has = true
					return 1, false
				}
			}
			// should never get here
			panic("Impossible lack of space for snake!")
		} else if len(direcs) == 1 {
			return direcs[0], true
		} else {
			return direcs[(*quickrand.RandInts(len(direcs)-1, 1, false))[0]], true
		}
	}
	for i := range ga.ind.snakes {
		// snakes should go in a straight line in the direction they are
		// facing and have a 1/3 chance for them to turn into a random valid direction.
		// when they hit the edge they should turn into a random valid direction
		direc := ga.Lawn[ga.ind.snakes[i]].Snake.Direc
		ok := true
		if inBounds(ga.ind.snakes[i], ga.ind.snakes[i]+inc_map[direc], direc) {
			dice := (*quickrand.RandInts(2, 1, false))[0]
			// roll for 1/3 chance to change direction to random new valid one
			if dice == 0 {
				direc, ok = getDirec(i, direc)
			}
			// ok is to make sure that snake was not tangled (b/c if so it shouldn't move forwards after the direc function sets it)
			if ok {
				ga.Lawn[ga.ind.snakes[i]].Snake.Has = false
				ga.ind.snakes[i] += inc_map[direc]
				ga.Lawn[ga.ind.snakes[i]].Snake.Direc = direc
				ga.Lawn[ga.ind.snakes[i]].Snake.Has = true
			}
		} else {
			direc, ok = getDirec(i, ga.Lawn[ga.ind.snakes[i]].Snake.Direc)
			if ok {
				ga.Lawn[ga.ind.snakes[i]].Snake.Has = false
				ga.ind.snakes[i] += inc_map[direc]
				ga.Lawn[ga.ind.snakes[i]].Snake.Direc = direc
				ga.Lawn[ga.ind.snakes[i]].Snake.Has = true
			}
		}
	}
}

// for rescaling offset values when window size changes
func (ga *Garden) ScaleOffset(len float32) {
	for i := 0; i < LawnLength; i++ {
		for j := 0; j < LawnLength; j++ {
			ga.Lawn[i*LawnLength+j].Off = f32.Point{float32(j) * (len / LawnLength), float32(i) * (len / LawnLength)}
		}
	}
}

// should be called after every gopher movement and snake movement.
// updates Garden state, for ex: if snake is on same square as gopher,
// gopher is dead and game is over.
// returns the 2 Grass inds which need to be updated
func (ga *Garden) Update() [2]int {
	// check if gopher is dead from snake
	var in_val_arg []int
	if in_val_arg = ga.ind.snakes[:]; quickrand.InVals(ga.ind.player, &in_val_arg) {
		ga.Dead = true
		// no suicide carrot eating. if gopher dies he will not collect the carrot on the square he died
		// check if gopher ate carrot and handle carrot regen
	} else if in_val_arg = ga.ind.carrots[:]; quickrand.InVals(ga.ind.player, &in_val_arg) {
		ga.Score++
		ga.Lawn[ga.ind.player].Carrot = false
		// find which carrot ind is same as player and set it to -1 so handleCarrots() knows it is not in the garden
		for i := 0; i < Carrots; i++ {
			if ga.ind.carrots[i] == ga.ind.player {
				ga.ind.carrots[i] = -1
				break
			}
		}
		ga.handleCarrots()
	}
	if ga.prev_hole != -1 {
		return [2]int{ga.ind.player, ga.prev_hole}
	}
	var reverse_inc = map[int]int{0: LawnLength, 1: -LawnLength, 2: -1, 3: 1}
	return [2]int{ga.ind.player, ga.ind.player + reverse_inc[ga.Lawn[ga.ind.player].Player.Direc]}
}
