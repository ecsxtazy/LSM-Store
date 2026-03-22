package memtable

import (
	"reflect"
	"testing"
)

func TestPutGet(t *testing.T) {
	mt := New()
	mt.Put("a", []byte("1"))
	mt.Put("b", []byte("2"))

	v, ok := mt.Get("a")
	if !ok || string(v) != "1" {
		t.Fatalf("unexpected get: ok=%v v=%q", ok, v)
	}
}

func TestOverwrite(t *testing.T) {
	mt := New()
	mt.Put("a", []byte("1"))
	mt.Put("a", []byte("2"))

	v, ok := mt.Get("a")
	if !ok || string(v) != "2" {
		t.Fatalf("unexpected get: ok=%v v=%q", ok, v)
	}
}

func TestDeleteTombstone(t *testing.T) {
	mt := New()
	mt.Put("a", []byte("1"))
	mt.Delete("a")

	if v, ok := mt.Get("a"); ok || v != nil {
		t.Fatalf("expected tombstone, got ok=%v v=%q", ok, v)
	}
}

func TestRangeOrder(t *testing.T) {
	mt := New()
	mt.Put("b", []byte("2"))
	mt.Put("a", []byte("1"))
	mt.Put("c", []byte("3"))

	var got []string
	mt.Range(func(rec Record) bool {
		if rec.Tombstone {
			return true
		}
		got = append(got, rec.Key)
		return true
	})

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order mismatch: got=%v want=%v", got, want)
	}
}

func TestRangeIncludesTombstone(t *testing.T) {
	mt := New()
	mt.Put("a", []byte("1"))
	mt.Delete("a")

	var got []Record
	mt.Range(func(rec Record) bool {
		got = append(got, rec)
		return true
	})

	if len(got) != 1 {
		t.Fatalf("expected one record, got %d", len(got))
	}
	if got[0].Key != "a" || !got[0].Tombstone {
		t.Fatalf("expected tombstone record for key a, got %+v", got[0])
	}
}
