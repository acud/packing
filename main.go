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

	// fmt.Println(fields)
	fields := [][]byte{c.A, c.b, c.c, c.d, c.e, c.f}
	// nn, fieldsLefttt := packFields(fields)
	root := packFields(fields)
	// fmt.Println("nn", root)

	return root
	// n := &node{}
	// rootNode := n
	// childCount := 0
	// n.a = make([]byte, 32)
	// n.b = make([]byte, 32)
	// n.c = make([]byte, 32)

	// fmt.Println("FIELDS", fields)
	// for i := 0; i < len(fields); i++ {
	// for _, node := range [][]byte{n.a, n.b, n.c} {
	// ctr := 0
	// fmt.Println("node", node, "ctr", ctr)
	// COPY:
	// ctr += copy(node[ctr:], fields[i])

	//fmt.Println("copy node", node, "ctr", ctr)
	//if i == len(fields)-1 {
	//break
	//}
	//if len(fields[i+1])+ctr < BYTES_PER_CHUNK {
	//fmt.Println("next field fits")
	//i++
	//goto COPY
	//}
	//i++
	//}
	//if i != len(fields)-1 {
	//fmt.Println("more fields, need new node")
	//fmt.Println("root", rootNode)
	//rootNode = &node{}
	//switch childCount {
	//case 0:
	//rootNode.a1 = n
	//case 1:
	//rootNode.b1 = n
	//case 2:
	//rootNode.c1 = n
	//}
	//childCount++
	//n.a = make([]byte, 32)
	//n.b = make([]byte, 32)
	//n.c = make([]byte, 32)
	//}
	//}

	// fmt.Println("n", n)

	// return *rootNode
}

//func packFields(fields [][]byte) (node, [][]byte) {
//fmt.Println("pack fields", fields)
//if len(fields) == 0 {
//return node{}, nil
//}
//n := node{}
//n.a = make([]byte, 32)
//n.b = make([]byte, 32)
//n.c = make([]byte, 32)

// i := 0
// for j := 0; j < 3; j++ {
// ctr := 0
// node := make([]byte, 32)
// fmt.Println("node", node, "ctr", ctr)
// COPY:
// ctr += copy(node[ctr:], fields[i])

//fmt.Println("copy node", node, "ctr", ctr)
//if i == len(fields)-1 {
//break
//}
//if len(fields[i+1])+ctr < BYTES_PER_CHUNK {
//fmt.Println("next field fits")
//i++
//goto COPY
//}
//i++
//n.leaf = append(n.leaf, node)
//}
//if len(n.leaf) == 3 {
//father := node{}
//}
//if i != len(fields)-1 {
//n1, fieldsLeft := packFields(fields[i:])
//n2, fieldsLeft := packFields(fieldsLeft)
//n3, fieldsLeft := packFields(fieldsLeft)
//n.a1 = &n1
//n.b1 = &n2
//n.c1 = &n3
//}
//return n, nil
//}

func packFields(fields [][]byte) *treeNode {
	fmt.Println("pack fields", fields)
	tree := &treeNode{}
	i := 0
	node := make([]byte, 32)

	for {
		ctr := 0
		fmt.Println("i", i, "node", node, "ctr", ctr)
	COPY:
		ctr += copy(node[ctr:], fields[i])
		fmt.Println("copy node", node, "ctr", ctr)
		if i == len(fields)-1 {
			ok, root := tree.insert(node)
			if !ok {
				tx := &treeNode{}
				tx.child = append(tx.child, tree)
				fmt.Println("new tree insert", tx)

				ook, newRoot := tx.insert(node)
				if !ook {
					panic("not ok")
				}
				tree = newRoot
				root = newRoot
			}
			// panic("never")
			return root
		}

		if len(fields[i+1])+ctr < BYTES_PER_CHUNK {
			fmt.Println("next field fits")
			i++
			goto COPY
		}
		fmt.Println("i", i)
		i++
		var ok bool
		ok, newRoot := tree.insert(node)

		if !ok {
			fmt.Println("not ok 1")

			tx := &treeNode{}
			tx.child = append(tx.child, tree)
			fmt.Println("new tree", tx)

			ook, newRoot := tx.insert(node)
			if !ook {
				panic("not ok")
			}
			tree = newRoot
		}
		tree = newRoot
		node = make([]byte, 32)
	}
	//if i != len(fields)-1 {
	//n1, fieldsLeft := packFields(fields[i:])
	//n2, fieldsLeft := packFields(fieldsLeft)
	//n3, fieldsLeft := packFields(fieldsLeft)
	//n.a1 = &n1
	//n.b1 = &n2
	//n.c1 = &n3
	//}
	return tree
}

type treeNode struct {
	child []*treeNode // if not nil then we're in an intermediate, otherwise leaf
	leaf  [][]byte
}

const deg = 3

func (t *treeNode) insert(chunk []byte) (bool, *treeNode) {
	if len(t.child) != 0 {
		// case 1 - we can still fit on the existing subtree
		for i, child := range t.child {
			fmt.Println("child try insert", i)
			ok, root := child.insert(chunk)
			if ok {
				fmt.Println("child inserted", i)
				return ok, root
			}
		}
		fmt.Println("tried all children")
		if len(t.child) < deg {
			fmt.Println("can still add child")
			node := &treeNode{}
			node.insert(chunk)
			t.child = append(t.child)
			return true, t
		} else {
			fmt.Println("create new father")
			// create new father,
			father := &treeNode{}
			father.child = append(father.child, t)
			return father.insert(chunk)
		}
	}
	if len(t.leaf) == deg {
		fmt.Println("cant add anymore")
		return false, t
	}
	fmt.Println("adding...")
	t.leaf = append(t.leaf, chunk)
	return true, t
}

func main() {
	fmt.Println("abcd")
	// cc := container{
	//a: uint32(1),
	//}
	//cc.pack()
}
