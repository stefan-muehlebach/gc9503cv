package main

import (
	"testing"
)

var (
	p0 = NewVector(0, 2, 0)
	p1 = NewVector(-1, 1, 0)
	p2 = NewVector(1, 1, 0)
	p3 = NewVector(-1, -1, 0)
	p4 = NewVector(1, -1, 0)
	p5 = NewVector(0, -2, 0)

	p10 = NewVector(0, 2, 1)
	p11 = NewVector(-1, 1, 1)
	p12 = NewVector(1, 1, 1)
	p13 = NewVector(-1, -1, 1)
	p14 = NewVector(1, -1, 1)
	p15 = NewVector(0, -2, 1)
)

func TestIntersects(t *testing.T) {
	if inter := Intersects(p1, p2, p3, p4); inter {
		t.Errorf("p1,p2  and  p3,p4  shouldn't intersect")
	}
	if inter := Intersects(p1, p4, p2, p3); !inter {
		t.Errorf("p1,p4  and  p2,p3  should intersect")
	}
	if inter := Intersects(p1, p4, p12, p13); inter {
		t.Errorf("p1,p4  and  p12,p13  shouldn't intersect")
	}
	if inter := Intersects(p11, p14, p12, p13); !inter {
		t.Errorf("p11,p14  and  p12,p13  should intersect")
	}
	if inter := Intersects(p1, p14, p2, p13); !inter {
		t.Errorf("p1,p14  and  p2,p13  should intersect")
	}
	if inter := Intersects(p0, p15, p1, p14); !inter {
		t.Errorf("p0,p15  and  p1,p14  should intersect")
	}
}
