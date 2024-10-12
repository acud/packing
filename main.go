package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

const BYTES_PER_CHUNK = 32

type (
	hashFn = hash.Hash
)

var (
	allZeros []byte
	zero     [32]byte
)

func init() {
	var v [32 * 3]byte
	allZeros = hashData(v[:])
}

func hashData(d []byte) []byte {
	v := sha256.Sum256(d)
	return v[:]
}

type node struct {
	a1, b1, c1 *node // if not nil then we're in an intermediate, otherwise leaf
	a, b, c    []byte
}

func (n node) hash() []byte {
	// we're in an intermediate
	if n.a1 != nil {

		a1H := n.a1.hash()
		fmt.Println("a1H", hex.EncodeToString(a1H))
		b1H := n.b1.hash()
		fmt.Println("b1H", hex.EncodeToString(b1H))
		c1H := n.c1.hash()
		fmt.Println("c1H", hex.EncodeToString(c1H))
		v := []byte{}
		v = append(v, a1H...)
		v = append(v, b1H...)
		v = append(v, c1H...)
		return hashData(v)
	}

	if n.a == nil && n.b == nil && n.c == nil {
		return allZeros
	}

	v := []byte{}

	if n.a == nil {
		fmt.Println("a zero")
		v = append(v, zero[:]...)
	} else {
		v = append(v, n.a...)
		addPad(v, BYTES_PER_CHUNK-len(n.a))
	}
	if n.b == nil {
		fmt.Println("b zero")
		v = append(v, zero[:]...)
	} else {
		v = append(v, n.b...)
		addPad(v, BYTES_PER_CHUNK-len(n.a))
	}
	if n.c == nil {
		fmt.Println("a zero")
		v = append(v, zero[:]...)
	} else {
		v = append(v, n.c...)
		addPad(v, BYTES_PER_CHUNK-len(n.a))
	}
	return hashData(v)
}

func addPad(s []byte, size int) {
	s = append(s, zero[:][:size]...)
}

type container struct {
	A []byte // 8
	b []byte // 16
	c []byte // 16
	d []byte // 32
	e []byte // 8
	f []byte // 8
}

func (c container) pack() *treeNode {
	// fields := reflect.VisibleFields(reflect.TypeOf(c))
	// for _, field := range fields {
	// vv := reflect.ValueOf(c)
	// vvr := reflect.Indirect(vv.FieldByName(field.Name)).Bytes()

	//fmt.Println("printing field", field.Name, vvr)
	//}

	fields := [][]byte{c.A, c.b, c.c, c.d, c.e, c.f}
	root := packFields(fields)

	return root
}

func packFields(fields [][]byte) *treeNode {
	fmt.Println("pack fields", fields)
	tree := &treeNode{child: []*treeNode{{}}}

	i := 0
	node := make([]byte, 32)

	for {
		ctr := 0
	COPY:
		ctr += copy(node[ctr:], fields[i])
		if i == len(fields)-1 {
			// last field. insert the node and return the new root
			ok, root := tree.insert(node)
			if !ok {
				panic(1)
			}
			return root
		}

		if len(fields[i+1])+ctr < BYTES_PER_CHUNK {
			// next field fits
			i++
			goto COPY
		}
		i++
		var ok bool
		ok, newRoot := tree.insert(node)

		if !ok {
			panic(1)
		}
		tree = newRoot
		node = make([]byte, 32)
	}
	return tree
}

type treeNode struct {
	child []*treeNode // if not nil then we're in an intermediate, otherwise leaf
	leaf  [][]byte
}

const deg = 3

func (t *treeNode) hash() []byte {
	if len(t.child) > 0 {
		// there's child nodes to call
		hashes := []byte{}
		for i, child := range t.child {
			fmt.Println("calling hash on child", i)
			hash := child.hash()
			fmt.Println("child hash", i, hex.EncodeToString(hash))
			hashes = append(hashes, hash...)
		}

		fmt.Println("adding zero hash", deg, len(t.child))
		for i := deg - len(t.child); i > 0; i-- {
			hashes = append(hashes, allZeros...)
		}
		return hashData(hashes)
	}

	v := []byte{}

	for _, item := range t.leaf {
		v = append(v, item...)
		addPad(v, BYTES_PER_CHUNK-len(item))
	}
	for i := deg - len(t.leaf); i > 0; i-- {
		v = append(v, zero[:]...)
	}

	return hashData(v)
}

func (t *treeNode) insert(chunk []byte) (bool, *treeNode) {
	if len(t.child) != 0 {
		// case 1 - we can still fit on the existing subtree
		for _, child := range t.child {
			ok, _ := child.insert(chunk)
			if ok {
				return ok, t
			}
		}
		if len(t.child) < deg {
			// can still add a child node here
			node := &treeNode{}
			node.insert(chunk)
			t.child = append(t.child, node)
			return true, t
		} else {
			// case 2 - we can't fit on the subtree
			// hence we create a new father, add this tree to it
			// and call insert on the new father which should always succeed
			father := &treeNode{}
			father.child = append(father.child, t)
			return father.insert(chunk)
		}
	}
	if len(t.leaf) == deg {
		// leaf node full
		return false, t
	}
	// leaf node can still fit a value
	t.leaf = append(t.leaf, chunk)
	return true, t
}

func main() {
	fmt.Println("abcd")
}
