package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func mustDecode(s string) []byte {
	result, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return result
}

func TestFixture(t *testing.T) {
	a := mustDecode("aaaaaaaaaaaaaaaa")
	b := mustDecode("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	c := mustDecode("cccccccccccccccccccccccccccccccc")
	d := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	e := mustDecode("eeeeeeeeeeeeeeee")
	f := mustDecode("ffffffffffffffff")

	v := container{}
	v.A = a
	v.BRandonName = b
	v.C = c
	v.D = d
	v.E = e
	v.F = f

	treeRoot := v.pack()
	result := treeRoot.hash()
	fmt.Println("got result hash", hex.EncodeToString(result))
	expRes, err := hex.DecodeString("4278118c38f02679efc01a9075510abe00747b01705c9add495053d88604ce95")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(result, expRes) != 0 {
		t.Log(hex.EncodeToString(result))
		t.Fatal("hash mismatch")
	}
}

func TestDeeper(t *testing.T) {
	a := mustDecode("aaaaaaaaaaaaaaaa")
	b := mustDecode("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	c := mustDecode("cccccccccccccccccccccccccccccccc")
	d := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	e := mustDecode("eeeeeeeeeeeeeeee")
	f := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	g := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	h := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	i := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	j := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	k := mustDecode("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")

	v := bigcontainer{}
	v.A = a
	v.BRandonName = b
	v.C = c
	v.D = d
	v.E = e
	v.F = f
	v.G = g
	v.H = h
	v.I = i
	v.J = j
	v.K = k

	treeRoot := v.pack()
	result := treeRoot.hash()
	spew.Dump(treeRoot)
	fmt.Println("got result hash", hex.EncodeToString(result))
}

type bigcontainer struct {
	A           []byte // 8
	BRandonName []byte // 16
	C           []byte // 16
	D           []byte // 32
	E           []byte //
	F           []byte //
	G           []byte //
	H           []byte //
	I           []byte //
	J           []byte //
	K           []byte //
}

func (c bigcontainer) pack() *treeNode {
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

/*
	{
	  a: 0xaaaaaaaaaaaaaaaa
	  b: 0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
	  c: 0xcccccccccccccccccccccccccccccccc
	  d: 0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd
	  e: 0xeeeeeeeeeeeeeeee
	  f: 0xffffffffffffffff
	}
	{
	  a: 0xaaaaaaaaaaaaaaaa
	  b: 0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
	  c: 0xcccccccccccccccccccccccccccccccc
	  d: 0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd
	  e: 0xeeeeeeeeeeeeeeee
	  f: 0xffffffffffffffff
	}

which hash to an intermediate layer of roots:

0x49a1a674384df6413ca11c56bae33b6bc7590c1bcd9e7225f46e615e1e1246d4
0xe45a4dfc1ffcb25e44beb6c497f251c658a3d0f4454a097d59f7ef69e580c730
0x2ea9ab9198d1638007400cd2c3bef1cc745b864b76011a0e1bc52180ac6452d4

SHA256.digest(concat([

	fromHex("aaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb0000000000000000"),
	fromHex("cccccccccccccccccccccccccccccccc00000000000000000000000000000000"),
	fromHex("dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"),

])).toHex()
// Returns 49a1a674384df6413ca11c56bae33b6bc7590c1bcd9e7225f46e615e1e1246d4
which hash to the tree root

0x4278118c38f02679efc01a9075510abe00747b01705c9add495053d88604ce95
*/
