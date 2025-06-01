package encoder

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/emirpasic/gods/v2/maps/treemap"
)

type rankNode struct {
	rank  int
	index int
	prev  *rankNode
	next  *rankNode
}

func newRankNode(rank int, index int, prev *rankNode) *rankNode {
	return &rankNode{
		rank:  rank,
		index: index,
		prev:  prev,
	}
}

func (r *rankNode) ToString() string {
	return fmt.Sprintf("rankNode{rank=%d, index=%d}", r.rank, r.index)
}

func RemoveNode(nodeMap *linkedhashmap.Map[int, *rankNode], rankMap *treemap.Map[int, *linkedhashmap.Map[int, *rankNode]], node *rankNode) {
	if nodeMap.Size() == 1 {
		if _, ok := nodeMap.Get(node.index); ok {
			rankMap.Remove(node.rank)
		} else {
			panic("nodeMap.containsKey(node.index);")
		}
	} else {
		nodeMap.Remove(node.index)
	}
}

func CalculateTokensLarge(tokenEncoder *TokenEncoder, maxTokenCount int, keepEncodings bool, out []int, match []byte) int {
	rankMap := treemap.New[int, *linkedhashmap.Map[int, *rankNode]]()

	var prev *rankNode = nil
	for i := 0; i < len(match)+1; i++ {
		rank := tokenEncoder.Encode(match, i, i+2)
		node := newRankNode(rank, i, prev)
		if prev != nil {
			prev.next = node
		}
		prev = node

		rankNodeMap, ok := rankMap.Get(rank)
		if !ok {
			rankNodeMap = linkedhashmap.New[int, *rankNode]()
			rankMap.Put(rank, rankNodeMap)
		}
		rankNodeMap.Put(i, node)
	}

	if _, ok := rankMap.Get(MAX_RANK); !ok {
		panic("rankMap.containsKey(MAX_RANK)")
	}

	tokenCount := len(match)
	for tokenCount > 2 && rankMap.Size() > 1 {
		firstIt := rankMap.Iterator()
		var minNodeMap *linkedhashmap.Map[int, *rankNode] = nil
		if firstIt.First() {
			minNodeMap = firstIt.Value()
			rankMap.Remove(firstIt.Key())
		}

		it := minNodeMap.Iterator()
		for it.Begin(); it.Next(); {
			minNode := it.Value()
			minRank := minNode.rank
			if minRank == MAX_RANK {
				panic("minRank == MAX_RANK")
			}

			previousNode := minNode.prev
			nextNode := minNode.next
			nextNextNode := nextNode.next
			nextNextNextNode := nextNextNode.next

			if previousNode != nil {
				newRank := tokenEncoder.Encode(match, previousNode.index, nextNextNode.index)
				if previousNode.rank != newRank {
					if previousNode.rank == minRank {
						panic("previousNode.rank == minRank")
					}
					if nodeMap, ok := rankMap.Get(previousNode.rank); ok {
						RemoveNode(nodeMap, rankMap, previousNode)
					}
					previousNode.rank = newRank
					nodeMap, ok := rankMap.Get(newRank)
					if !ok {
						nodeMap = linkedhashmap.New[int, *rankNode]()
						rankMap.Put(newRank, nodeMap)
					}
					nodeMap.Put(previousNode.index, previousNode)
				}
			}

			newRankIndex := math.MaxInt32
			if nextNextNextNode != nil {
				newRankIndex = nextNextNextNode.index
			}
			newRank := tokenEncoder.Encode(match, minNode.index, newRankIndex)
			minNode.rank = newRank
			nodeMap, ok := rankMap.Get(newRank)
			if !ok {
				nodeMap = linkedhashmap.New[int, *rankNode]()
				rankMap.Put(newRank, nodeMap)
			}
			nodeMap.Put(minNode.index, minNode)

			minNode.next = nextNextNode
			nextNextNode.prev = minNode
			if nextNode.rank != MAX_RANK {
				if nextNode.rank != minRank {
					if nodeMap, ok := rankMap.Get(nextNode.rank); ok {
						RemoveNode(nodeMap, rankMap, nextNode)
					}
				} else {
					it.Next()
				}
			}

			tokenCount--
		}
	}

	if keepEncodings {
		headNodeMap, _ := rankMap.Get(MAX_RANK)
		for head, _ := headNodeMap.Get(0); head.next != nil && len(out) < maxTokenCount; head = head.next {
			token := tokenEncoder.Encode(match, head.index, head.next.index)
			if token == MAX_RANK {
				panic("Token should not be MAX_RANK")
			}
			out = append(out, token)
		}
	}

	return tokenCount
}
