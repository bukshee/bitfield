package bitfield

import (
	"testing"
)

func Test1(t *testing.T) {
	if New(0) != nil || New(-2) != nil {
		t.Error("should be nil")
	}
	if New(65).SetAll().OnesCount() != 65 {
		t.Error("should be 65")
	}
	if New(3).Equal(New(4)) {
		t.Error("should be false")
	}
	if New(3).Set(0).Set(-1).Not().OnesCount() != 1 {
		t.Error("should be 1")
	}
	bf := New(129).Set(0).Set(-1).Clear(123).Clear(-3).Not().Not()
	if bf.OnesCount() != 2 {
		t.Error("should be 2")
	}
	if !bf.Get(0) || !bf.Get(-bf.Len()) || !bf.Get(bf.Len()) {
		t.Error("should be true")
	}
	if bf.And(New(129).Set(0).Set(1)).OnesCount() != 1 {
		t.Error("should be 1")
	}

	if bf.And(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf.Or(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}
	if bf.Xor(New(121)).OnesCount() != 1 {
		t.Error("should be 1")
	}

	bf.Set(73).Set(-2).ClearAll().Set(-1)
	if !bf.Equal(New(129).Set(128)) {
		t.Error("should be equal")
	}
	if bf.Equal(New(129).Not()) {
		t.Error("should be not equal")
	}
	if bf.Get(127) {
		t.Error("should be false")
	}
	if !bf.Get(-1) {
		t.Error("should be true")
	}

	if !New(4).Flip(-1).Equal(New(4).Set(-1)) {
		t.Error("should be equal")
	}
	if !New(4).Flip(-1).Flip(-1).Equal(New(4)) {
		t.Error("should be equal")
	}

	bf2 := bf.Copy()
	if !bf.Equal(bf2) || bf.OnesCount() != bf2.OnesCount() || bf.Len() != bf2.Len() {
		t.Error("should be equal")
	}
	if bf2.Xor(bf2).OnesCount() != 0 {
		t.Error("should be 0")
	}

	bf2 = bf.Copy()
	if !bf2.Set(11).Or(bf).Get(11) {
		t.Error("should be true")
	}

}