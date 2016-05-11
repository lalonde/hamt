package hamt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestEmptyTrie(t *testing.T) {
	trie := New()

	_,e := trie.Get([]byte("EMPTY"))

	if e==nil {
		t.Errorf("Not finding a key, which we wont find in an empty trie, must return an err")
	}
}

func TestInsertIntoTrie(t *testing.T) {
	key1 := []byte("store key1")
	key2 := []byte("store key2")

	v1 := "value 1"
	v2 := "value 2"


	root := New()

	root.Insert( key1, v1 )
	root.Insert( key2, v2 )

	vg1, e1 := root.Get( key1 )
	vg2, e2 := root.Get( key2 )

	if vg1 != v1 || vg2 != v2 {
		t.Errorf("Set values for keys[%v=%v,%v=%v] do not match returned [%v=%v,%v=%v]",key1,v1,key2,v2,key1,vg1,key2,vg2)
		t.Errorf("error return: %v, %v", e1, e2)
	}

	intKeys := make([][]byte, 256)
	for i := 0; i < 256; i++ {
		intKeys[i] = IntKey(i)
	}

	stringKeys := make([][]byte, 256)
	for i := 0; i < 256; i++ {
		stringKeys[i] = StringKey(fmt.Sprint("String key", i))
	}

	insertAndAssureStorageForKeys(intKeys, t)
	insertAndAssureStorageForKeys(stringKeys, t)
}

func insertAndAssureStorageForKeys(keys [][]byte, t *testing.T) {
	tree := New()
	makeValueForKey := func (key interface{}) string { return fmt.Sprintf("Value for %v", key) }

	for _, k := range keys {
		v := makeValueForKey(k)
		tree.Insert(k, v)
		gv, ge := tree.Get(k)
		if ge != nil || gv != v {
			t.Errorf("We blewit! %v not equal to %v for %v", v, gv, k)
		}
	}

	for _, k := range keys {
		expectedValue := makeValueForKey(k)
		gv, ge := tree.Get(k)
		if ge != nil || gv != expectedValue {
			t.Errorf("After full insert of keys we got inequality for key %v.  Expected %v but got %v", k, expectedValue, gv)
		}
	}
}


func BenchmarkGoMapIntInsert(b *testing.B) {
	m := make(map[int] int)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}

func BenchmarkGoMapStringInsert(b *testing.B) {
	m := make(map[string] string)
	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("String Key %v", i)
		v := fmt.Sprintf("String Val %v", i)
		m[k] = v
	}
}


func BenchmarkHamtIntInsert(b *testing.B) {
	t := New()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, i)
		t.Insert(buf.Bytes(), i)
	}
}

func BenchmarkHamtStringInsert(b *testing.B) {
	t := New()
	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("String Key %v", i)
		v := fmt.Sprintf("String Val %v", i)
		t.Insert([]byte(k), v)
	}
}
