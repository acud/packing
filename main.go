package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"reflect"
	"sync"
)

const (
	BYTES_PER_CHUNK = 32
	DEGREE          = 3
)

var (
	allZeros []byte
	zero     [32]byte
)

func init() {
	var v [32 * 3]byte
	res := sha256.Sum256(v[:])
	allZeros = res[:]
}

type treeNode struct {
	child []*treeNode // if not nil then we're in an intermediate, otherwise leaf
	leaf  [][]byte
}

var pool = &sync.Pool{
	New: func() any {
		return sha256.New()
	},
}

func getHasher() hash.Hash {
	return pool.Get().(hash.Hash)
}

func putHasher(h hash.Hash) {
	h.Reset()
	pool.Put(h)
}

func (t *treeNode) hash() []byte {
	if len(t.child) > 0 {
		// there're child nodes to call
		hashes := getHasher()
		defer putHasher(hashes)

		for _, child := range t.child {
			_, err := hashes.Write(child.hash())
			if err != nil {
				panic(err)
			}
		}

		for i := DEGREE - len(t.child); i > 0; i-- {
			// append the subtree hash of 3 nodes of all zero leaves
			// when they don't exist so we don't need to allocate the
			// subtree
			hashes.Write(allZeros)
		}
		return hashes.Sum(nil)
	}

	hasher := getHasher()
	defer putHasher(hasher)

	for _, item := range t.leaf {
		_, err := hasher.Write(item)
		if err != nil {
			panic(err)
		}
		hasher.Write(zero[:BYTES_PER_CHUNK-len(item)])
	}
	for i := DEGREE - len(t.leaf); i > 0; i-- {
		// append the zero hash of sha256(32 zeros) so we don't
		// need to allocate it whenever there's no leaf node
		hasher.Write(zero[:])
	}

	return hasher.Sum(nil)
}

func (t *treeNode) insert(chunk []byte) (bool, *treeNode) {
	// fmt.Println("insert")
	if len(t.child) != 0 {
		// case 1 - we can still fit on the existing subtree
		for _, child := range t.child {
			ok, _ := child.insert(chunk)
			if ok {
				return ok, t
			}
		}
		if len(t.child) < DEGREE {
			// can still add a child node here
			node := &treeNode{}
			node.insert(chunk)
			t.child = append(t.child, node)
			return true, t
		} else {
			// case 2 - we can't fit on the subtree
			// hence we create a new father, add this tree to it
			// and call insert on the new father which should always succeed
			ff := &treeNode{}
			ff.insert(chunk)
			father := &treeNode{
				child: []*treeNode{t, ff},
			}
			return true, father
		}
	}
	if len(t.leaf) == DEGREE {
		// leaf node full
		return false, t
	}
	// leaf node can still fit a value
	t.leaf = append(t.leaf, chunk)
	return true, t
}

type container struct {
	A           []byte // 8
	BRandonName []byte // 16
	C           []byte // 16
	D           []byte // 32
	E           []byte // 8
	F           []byte // 8
}

func (c container) pack() *treeNode {
	var fieldsToPack [][]byte
	fields := reflect.VisibleFields(reflect.TypeOf(c))
	containerVal := reflect.ValueOf(c)
	for _, field := range fields {
		vvr := reflect.Indirect(containerVal.FieldByName(field.Name)).Bytes()
		fieldsToPack = append(fieldsToPack, vvr)
	}
	root := packFields(fieldsToPack)
	return root
}

func packFields(fields [][]byte) *treeNode {
	c1 := &treeNode{}
	tree := &treeNode{child: []*treeNode{c1}}

	i := 0
	node := make([]byte, 32)

	for {
		ctr := 0
	COPY:
		ctr += copy(node[ctr:], fields[i])
		fmt.Println("copy", node, ctr, i)
		if i == len(fields)-1 {
			// last field. insert the node and return the new root
			ok, root := tree.insert(node)
			if !ok {
				panic("shouldnt happen")
			}
			return root
		}

		if len(fields[i+1])+ctr < BYTES_PER_CHUNK {
			// next field fits
			i++
			goto COPY
		}
		i++
		ok, newRoot := tree.insert(node)
		if !ok {
			panic("shouldnt happen")
		}
		tree = newRoot
		node = make([]byte, 32)
	}
	return tree
}

func main() {
	fmt.Println("abcd")
}
