package hamt

import (
	"testing"
)

func TestEmptyTrie(t *testing.T) {
	trie := New()
	root := trie.root
	if root.childBitmap > 0 {
		t.Errorf("This new node should be empty ")
	}
}

func TestSerarchEmpthTrie(t *testing.T) {
	root := New()

	_,_, err := root.Search([]byte("NOTFOUNDKEY"))

	if err == nil {
		t.Errorf("Search for a key that is not in trie needs to return not found err")
	}
}


func TestInsertIntoTrie(t *testing.T) {
	key1 := "store key1"
	key2 := "store key2"

	v1 := "value 1"
	v2 := "value 2"

	root := New()
	t.Logf("%v", (1 | 64 | 4))
	t.Logf("mask: %v", mask(hash([]byte(key1)), 0))
	t.Logf("bitpos: %v", bitpos(hash([]byte(key1)), 0))
	s0, p0, e0 := root.Search( []byte(key1) )
	t.Logf("Search [%v, %v, %v]", s0, p0, e0)

	n1 := root.Insert( []byte(key1), v1 )

	n2 := root.Insert( []byte(key2), v2 )

	s1, p1, e1 := root.Search( []byte(key1) )
	t.Logf("Search [%v, %v, %v]", s1, p1, e1)
	s2, p2, e2 := root.Search( []byte(key2) )
	t.Logf("Search [%v, %v, %v]", s2, p2, e2)
//	t.Logf("I stored %s for key %s and searching for key %s i got %s", n1.value, key1, key1, s1.value)

	if n1 != s1 {
		t.Errorf("kv 1 did not store and retrieve stored: %v searched: %v", n1, s1)
	}
	
	if n2 != s2 {
		t.Errorf("kv 2 did not store and retrieve stored: %v searched: %v", n2, s2)
	}
}