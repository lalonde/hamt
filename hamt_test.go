package hamt

import (
	"testing"
)

func TestEmptyTrie(t *testing.T) {
	trie := New()
	
	_,e := trie.Get("EMPTY")
	
	if e==nil {
		t.Errorf("Not finding a key, which we wont find in an empty trie, must return an err")
	}
}

func TestInsertIntoTrie(t *testing.T) {
	key1 := "store key1"
	key2 := "store key2"

	v1 := "value 1"
	v2 := "value 2"

	t.Logf("Hash of %v = %v", v1, hash(v1))
	t.Logf("Hash of %v = %v", v2, hash(v2))

	root := New()
	t.Logf("%v", (1 | 64 | 4))
	t.Logf("mask: %v", mask(hash(key1), 0))
	t.Logf("bitpos: %v", bitpos(hash(key1), 0))
	v0, e0 := root.Get( key1 )
	t.Logf("Get [%v, %v]", v0, e0)

	n1 := root.Insert( key1, v1 )
	t.Logf("%v", n1)
	n2 := root.Insert( key2, v2 )
	t.Logf("%v", n2)

	vg1, e1 := root.Get( key1 )
	t.Logf("Get [%v, %v]", vg1, e1)

	vg2, e2 := root.Get( key2 )
	t.Logf("Get [%v, %v]", vg2, e2)
//	t.Logf("I stored %s for key %s and searching for key %s i got %s", n1.value, key1, key1, s1.value)


}