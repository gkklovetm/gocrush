package gocrush

import (
	//"log"
	"math"
)

type StrawSelector struct {
	Straws map[Node]int64
}

func NewStrawSelector(n Node) *StrawSelector {
	var s = new(StrawSelector)
	s.Straws = make(map[Node]int64)
	if !n.IsLeaf() {
		var sortedNodes = n.GetChildren()
		var numLeft = len(sortedNodes)
		var straw float64 = 1.0
		var wbelow float64 = 0.0
		var lastw float64 = 0.0
		var i = 0
		for i < len(sortedNodes) {
			var current = sortedNodes[i]
			if current.GetWeight() == 0 {
				s.Straws[current] = 0
				i += 1
				continue
			}
			s.Straws[current] = int64(straw * 0x10000)
			i += 1
			if i == len(sortedNodes) {
				break
			}

			current = sortedNodes[i]
			var previous = sortedNodes[i-1]
			if current.GetWeight() == previous.GetWeight() {
				continue
			}
			wbelow += (float64(previous.GetWeight()) - lastw) * float64(numLeft)
			for j := 0; j < len(sortedNodes); j++ {
				if sortedNodes[j].GetWeight() == current.GetWeight() {
					numLeft -= 1
				} else {
					break
				}
			}
			var wnext float64 = float64(int64(numLeft) * (current.GetWeight() - previous.GetWeight()))
			var pbelow = wbelow / (wbelow + wnext)
			straw *= math.Pow(1.0/pbelow, 1.0/float64(numLeft))
			lastw = float64(previous.GetWeight())
		}
	}
	return s
}

func (s *StrawSelector) Select(input int64, round int64) Node {
	var result Node
	var hiScore = int64(-1)
	for child, straw := range s.Straws {
		var score = weightedScore(child, straw, input, round)
		if score > hiScore {
			result = child
			hiScore = score
		}
	}
	if result == nil {
		panic("Illegal state")
	}
	return result
}

func weightedScore(child Node, straw int64, input int64, round int64) int64 {

	var hash = hash3(input, Btoi(digestString(child.GetId())), round)
	hash = hash & 0xFFFF
	var weightedScore = hash * straw
	return weightedScore
}
