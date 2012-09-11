package hamt

import (
	"bytes"
	"errors"
	"hash/fnv"
)

const (
	fanoutLog2 = 6
	fanout uint = 1 << fanoutLog2
	fanMask uint = fanout -1
	maxDepth = 60/fanoutLog2
	keyNotFound = "Key not found"
)

type Trie struct {
	root *Node
	collision map[uint]interface{}
}

type Node struct {
	childBitmap uint64
	key []byte
	value interface{}
	children []*Node
	depth uint
}

//Shift key hash until leaf with matching key is found or key is not found
func (t *Trie) Search(key []byte) (node *Node, parent *Node, err error) {
	//Hash our key and look for it in the root
	hash := hash(key)
	parent = t.root
	node = t.root.children[mask(hash, 0)]
	
	//If our value is in the root or the hash does not map to root slot get out fast
	switch {
	case node == nil:
		return node, parent, errors.New(keyNotFound)
	case node.childBitmap == 0:
		if bytes.Equal(node.key, key) {
			return node, nil, nil
		}
	}

	//Got here means our slot is sub trie
	for depth := 1; depth < maxDepth; depth++ {
		shift := uint(depth * fanoutLog2)
		pos := bitpos(hash, shift)
		parent = node
		if cMap := node.childBitmap; (pos & cMap)  == 0 { //nothing in slot, not found
			return nil, node, errors.New(keyNotFound)
		} else {
			index := node.index(pos)			
			node = parent.children[index]
			//Is this a value or map
			if node.isLeaf() {
				// Is this who we are looking for?
				if bytes.Equal(node.key, key) {
					return node, parent, nil 
				}
			} else {
				continue
			}
		}
	}
	
	return nil, parent, errors.New(keyNotFound)
}

func (t *Trie) Insert(key []byte, value interface{}) (node *Node) {
	//TODO: miles to go.....
	hash := hash(key)

	snode, parent, _ := t.Search(key)

	if snode == nil {

		node = NewSub(key, parent.depth +1)
		node.value = value
		pos := bitpos(hash, parent.depth * fanoutLog2)
		parent.childBitmap |= pos
		if parent.depth == 0 {
			parent.children[mask(hash, 0)] = node
		} else {
			parent.children[parent.index(pos)] = node
		}
	}
	return node
}



func New() *Trie {
	return &Trie {
		root: NewSub(nil, 0),
		collision: make(map[uint]interface{}),
	}
}

func NewSub(key []byte, depth uint) *Node {
	return &Node {
		childBitmap: 0,
		key: key,
		value: nil,
		children: make([]*Node, fanout),
		depth: depth,
	}
}

func hash(b []byte) uint64 {
	h := fnv.New64()
	h.Write(b)
	return h.Sum64()
}

func shift(hash uint64, shift uint) uint64 {
	if shift == 0 {
		return hash
	}
	return hash >> shift
}

func mask(hash uint64, bshift uint) uint {
	return uint(shift(hash, bshift) & uint64(fanMask))
}

func bitpos(hash uint64, bshift uint) uint64 {
	return 1 << mask(hash, bshift)
}

func (n *Node) index(bit uint64) uint {
	return popcount_2(n.childBitmap & (bit - 1))
}

func (n *Node) isLeaf() bool {
	if n.key != nil && n.childBitmap == 0 {
		return true
	} 
	return false
}
