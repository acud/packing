package main

import (
	"encoding/hex"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

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
	v.b = b
	v.c = c
	v.d = d
	v.e = e
	v.f = f
	// copy(v.A[:], a)
	// copy(v.b[:], b)
	// copy(v.c[:], c)
	// copy(v.d[:], d)
	// copy(v.e[:], e)
	// copy(v.f[:], f)

	n := v.pack()
	spew.Dump(n)
	// result := n.hash()
	//expRes, err := hex.DecodeString("4278118c38f02679efc01a9075510abe00747b01705c9add495053d88604ce95")
	//if err != nil {
	//panic(err)
	//}
	//if bytes.Compare(result, expRes) != 0 {
	//t.Log(hex.EncodeToString(result))
	//panic("failed")
	//}
	//fmt.Println(v)
	//fmt.Println("node")
	//fmt.Println(n.a)
}
