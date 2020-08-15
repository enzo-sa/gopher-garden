package ui

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// reads and writes highscores data to and from highscores.txt

type hs struct {
	name  string
	score int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (u *Ui) updateHs() {
	// checks name and highscore in file and updates highscore or
	// adds new entry if name is unrecognized
	newEntry := true
	f, err := ioutil.ReadFile("highscores.txt")
	check(err)
	lines := strings.Split(string(f), "\n")
	for i, line := range lines {
		pair := strings.Split(line, " ")
		if len(pair) > 2 {
			tempName := ""
			for i := 0; i < len(pair)-2; i++ {
				tempName += pair[i]
			}
			pair[0] = tempName
			pair[1] = pair[len(pair)-1]
			pair = pair[0:2]
		}
		if pair[0] == u.name {
			newEntry = false
			// make sure that there is no trailing CR char for windows systems
			if strings.HasSuffix(pair[1], "\r") {
				pair[1] = pair[1][:len(pair[1])-1]
			}
			prevScore, err := strconv.Atoi(pair[1])
			check(err)
			if u.ga.Score > prevScore {
				// set new highscore
				pair[1] = strconv.Itoa(u.ga.Score)
				lines[i] = strings.Join(pair, " ")
			} else {
				// no need to rewrite
				return
			}
		}
	}
	// append new line with new name and highscore if name not found in file, i.e. newEntry is true
	if newEntry {
		lines = append(lines, fmt.Sprintf("%s %s", u.name, strconv.Itoa(u.ga.Score)))
	}
	rewrite := strings.Join(lines, "\n")
	err = ioutil.WriteFile("highscores.txt", []byte(rewrite), 0644)
	check(err)
}

// returns top 5 highscore and name value pairs
func (u *Ui) readHs() [5]hs {
	var top = [...]hs{hs{"", -1}, hs{"", -1}, hs{"", -1}, hs{"", -1}, hs{"", -1}}
	f, err := ioutil.ReadFile("highscores.txt")
	check(err)
	lines := strings.Split(string(f), "\n")
	if lines[0] == "" {
		return top
	}
	for _, line := range lines {
		pair := strings.Split(line, " ")
		if len(pair) > 1 {
			tempName := ""
			if len(pair) > 2 {
				for i := 0; i < len(pair)-1; i++ {
					tempName += pair[i]
				}
				pair[0] = tempName
				pair[1] = pair[len(pair)-1]
			}
			if strings.HasSuffix(pair[1], "\r") {
				pair[1] = pair[1][:len(pair[1])-1]
			}
			score, err := strconv.Atoi(pair[1])
			check(err)
			user := hs{pair[0], score}
			// insert score into top 5 position
			for i := 0; i < 5; i++ {
				if user.score > top[i].score {
					for j, temp := i, top[i]; j < 4; j++ {
						temp2 := top[j+1]
						top[j+1] = temp
						temp = temp2
					}
					top[i] = user
					break
				}
			}
		}
	}
	return top
}
