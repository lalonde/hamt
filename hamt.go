package hamt

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
)

const (
	fanoutLog2 = 6
	fanout uint = 1 << fanoutLog2
	fanMask uint = fanout -1
	maxDepth = 60/fanoutLog2
	keyNotFound = "Key not found"
)

type node interface {
	assoc(shift int, hash uint64, key interface{}, value interface{}) node
	without(shift int, hash uint64, key interface{}) node
	find(shift int, hash uint64, key interface{}) (value interface{}, err error)
	depth() int
}

type emptyNode struct {}

type PersistentMap struct {
	root node
//	collision map[uint]interface{}
}

type valueNode struct {
	key interface{}
	hash uint64
	value interface{}
	shift int
}

type bitmapNode struct {
	childBitmap uint64
	children []*node
	shift int
}

func (n *emptyNode) assoc(shift int, hash uint64, key interface{}, val interface{}) node {
	return &valueNode {key: key, hash: hash, value: val, shift: shift}
}
func (n *emptyNode) without(shift int, hash uint64, key interface{}) node {
	return n
}
func (n *emptyNode) find(shift int, hash uint64, key interface{}) (interface{}, error) {
	return nil, errors.New(keyNotFound)
}
func (n *emptyNode) depth() int {
	return 0
}

func (n *valueNode) assoc(shift int, hash uint64, key interface{}, val interface{}) node {
	return &valueNode {key: key, hash: hash, value: val, shift: shift}
}
func (n *valueNode) without(shift int, hash uint64, key interface{}) node {
	return n
}
func (n *valueNode) find(shift int, hash uint64, key interface{}) (value interface{}, err error) {
	if hash == n.hash {
		value = n.value
	} else {
		err = errors.New(keyNotFound)
	}
	return value, err
}
func (n *valueNode) depth() int {
	return n.shift
}

//Shift key hash until leaf with matching key is found or key is not found
func (t *PersistentMap) Get(key interface{}) (value interface{}, err error) {
	//Hash our key and look for it in the root
	hash := hash(key)
	value, err = t.root.find(0, hash, key)
	return
}

//	node = t.root.children[mask(hash, 0)]
//	
//	//If our value is in the root or the hash does not map to root slot get out fast
//	switch {
//	case node == nil:
//		return node, parent, errors.New(keyNotFound)
//	case node.childBitmap == 0:
//		if bytes.Equal(node.key, key) {
//			return node, nil, nil
//		}
//	}
//
//	//Got here means our slot is sub trie
//	for depth := 1; depth < maxDepth; depth++ {
//		shift := uint(depth * fanoutLog2)
//		pos := bitpos(hash, shift)
//		parent = node
//		if cMap := node.childBitmap; (pos & cMap)  == 0 { //nothing in slot, not found
//			return nil, node, errors.New(keyNotFound)
//		} else {
//			index := node.index(pos)			
//			node = parent.children[index]
//			//Is this a value or map
//			if node.isLeaf() {
//				// Is this who we are looking for?
//				if bytes.Equal(node.key, key) {
//					return node, parent, nil 
//				}
//			} else {
//				continue
//			}
//		}
//	}
//	
//	return nil, parent, errors.New(keyNotFound)
//}

func (t *PersistentMap) Insert(key interface{}, value interface{}) (n node) {
	//TODO: miles to go.....
		
	hash := hash(key)
	n = t.root.assoc(0, hash, key, value)
	if n.depth() == 0 {
		t.root = n
	}
	return n
}

//	snode, parent, _ := t.Search(key)
//
//	if snode == nil {
//
//		node = NewSub(key, parent.depth +1)
//		node.value = value
//		pos := bitpos(hash, parent.depth * fanoutLog2)
//		parent.childBitmap |= pos
//		if parent.depth == 0 {
//			parent.children[mask(hash, 0)] = node
//		} else {
//			parent.children[parent.index(pos)] = node
//
//	}
//	}
//	return node
//}



func New() *PersistentMap {
	return &PersistentMap {
		root: &emptyNode{},
	}
}

func hash(a interface{}) uint64 {
	buf := new(bytes.Buffer)
	h := fnv.New64()
	_, e := fmt.Fprint(buf, a)

	if e != nil {
		//doSomethingWithErrors?
	}

	h.Write(buf.Bytes())

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

func (n *bitmapNode) index(bit uint64) uint {
	return popcount_2(n.childBitmap & (bit - 1))
}
