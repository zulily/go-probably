package probably

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sort"
	"testing"
)

type ItemCounts []ItemCount

// Len is part of sort.Interface.
func (ic ItemCounts) Len() int {
	return len(ic)
}

// Swap is part of sort.Interface.
func (ic ItemCounts) Swap(i, j int) {
	ic[i], ic[j] = ic[j], ic[i]
}

// Less is part of sort.Interface
func (ic ItemCounts) Less(i, j int) bool {
	switch {
	case ic[i].Count < ic[j].Count:
		return true
	case ic[i].Count == ic[j].Count:
		return ic[i].Key < ic[j].Key
	default:
		return false
	}
}

func TestStreamGob(t *testing.T) {
	s := NewStreamTop(8, 3, 3)
	s2 := NewStreamTop(8, 3, 3)

	exp := []struct {
		s     string
		count int
	}{
		{"hello", 3},
		{"there", 2},
		{"world", 1},
	}

	for _, e := range exp {
		s.AddCount(e.s, uint32(e.count))
	}

	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(s); err != nil {
		t.Error(err)
	}
	if err := gob.NewDecoder(&buf).Decode(&s2); err != nil {
		t.Error(err)
	}

	sTop := s.GetTop()
	s2Top := s2.GetTop()

	sort.Sort(ItemCounts(sTop))
	sort.Sort(ItemCounts(s2Top))

	for i, v := range sTop {
		if v != s2Top[i] {
			t.Errorf("counts differ for %q", v)
		}
	}

	if !reflect.DeepEqual(sTop, s2Top) {
		t.Error("unmarshaled structure differs")
	}
}
