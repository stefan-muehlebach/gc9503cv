package main

import "math"

//-----------------------------------------------------------------------

// Die folgenden Datenstrukturen sind fuer die Offline-Speicherung von
// 3-dimensionalen Dingen gedacht. Sie stellen die Daten so dar, dass
// ein Benutzer mittels Texteditor damit arbeiten kann.
type ObjectData struct {
	// Enthaelt die Koordinaten aller Punkte, welche fuer dieses Objekt
	// von Relevanz sind. Die nachfolgenden Angaben ('Faces') beziehen
	// sich immer auf den Index des jeweiligen Punktes innerhalb dieses
	// Slices.
	Points []Vector
	// Liste von 'einfachen Flaechen' - das sind Flaechen, die keine
	// Aussparungen haben und die nicht mit weiteren Flaechen verbunden
	// werden sollen.
	// Jede Flaeche wird durch einen Polygonzug definiert, dessen Punkte
	// im gegenuhrzeigersinn angegeben werden muessen. Der Polygonzug wird
	// automatisch geschlossen und ist auch gleich dem Rand dieser Flaeche.
	// Beim Laden wird geprueft, ob die Punkte alle in einer Ebene liegen
	// und die Flaeche wird autom. in Dreiecke aufgeteilt.
	// TO DO: der Algorithmus zur Triangulation ist noch nicht einwandfrei.
	SimpleFaces [][]int
	// Liste von 'triangulierten Flaechen', das sind Flaechen, die aus
	// lauter Dreiecken zusammengesetzt sind. Damit lassen sich beliebige
	// Formen bauen - auch solche mit Aussparungen.
	TrigFaces []FaceData
}

// Die Daten zu einer triangulierten Flaeche, bei welcher der Benutzer die
// volle Kontrolle ueber die Triangulation hat.
type FaceData struct {
	// Mit 'Edges' koennen die einfassenden Begrenzungslinien (egal ob
	// Rand oder Aussparung) spezifiziert werden. Diese Polygonzuege
	// werden NICHT geschlossen.
	Edges [][]int
	// 'Trigs' ist eine Liste von Dreiecken, aus denen die Flaeche zusammen-
	// gesetzt ist. Auch hier gilt, dass die Punkte im Gegenuhrzeigersinn
	// aufgefuert werden muessen.
	Trigs []TrigData
}

// Mit diesem Datentyp schliesslich werden die einzelnen Dreiecke definiert.
// Damit laesst sich auch verhindern, dass aus Versehen laengere oder kuerzere
// Slices angegeben werden.
type TrigData [3]int

//----------------------------------------------------------------------------

var (
	// Die platonischen Koerper (5 Stueck) machen den Anfang der
	// Objektsammlung. Als plat. Koerper gilt jeder konvexe Polyeder, der
	//   a) nur durch kongruente, konvexe, regelmaessige Polygone begrenzt,
	//   b) dessen Seiten sich nicht schneiden (ausser an den Raendern)
	//   c) bei allen Eckpunkten die gleiche Anzahl Flaechen zusammentreffen
	//
	// Die folgende Tabelle zeigt die Eigenschaften der plat. Koerper:
	//
	//  Polyeder     | Ecken | Kanten | Flaechen
	// --------------+-------+--------+----------
	//  Tetraeder    |     4 |      6 |        4
	//  Hexaeder     |     8 |     12 |        6
	//  Oktaeder     |     6 |     12 |        8
	//  Dodekaeder   |    20 |     30 |       12
	//  Ikosaeder    |    12 |     30 |       20
	//
	// Als Koordinaten der Eckpunkte werden folgende Daten verwendet:
	// ('phi' ist das Verhaeltnis des goldenen Schnittes (1+sqrt(5))/2)
	//
	// Tetraeder : (1,1,1), (1,-1,-1), (-1,1,-1), (-1,-1,1)
	// Hexaeder  : (±1, ±1, ±1)
	// Oktaeder  : (±1, 0, 0), (0, ±1, 0), (0, 0, ±1)
	// Dodekaeder: (±1, ±1, ±1), (0, ±1/phi, ±phi), (±1/phi, ±phi, 0),
	//             (±phi, 0, ±1/phi)
	// Ikosaeder : (0, ±1, ±phi), (±1, ±phi, 0), (±phi, 0, ±1)

	// Constant for the golden ratio.
	phi  = math.Phi
	phiI = 1.0 / phi

	// Tetraeder - besteht aus 4 Dreiecken.
	Tetraeder = ObjectData{
		Points: []Vector{
			Vector{+1, +1, +1},
			Vector{+1, -1, -1},
			Vector{-1, +1, -1},
			Vector{-1, -1, +1},
		},
		SimpleFaces: [][]int{
			{0, 1, 2},
			{0, 3, 1},
			{0, 2, 3},
			{1, 3, 2},
		},
	}

	// Hexaeder (oder einfach Wuerfel, gell...) - besteht aus 6 Quadraten.
	Hexaeder = ObjectData{
		Points: []Vector{
			Vector{-1, -1, -1},
			Vector{-1, +1, -1},
			Vector{+1, +1, -1},
			Vector{+1, -1, -1},
			Vector{-1, -1, +1},
			Vector{-1, +1, +1},
			Vector{+1, +1, +1},
			Vector{+1, -1, +1},
		},
		SimpleFaces: [][]int{
			{0, 1, 2, 3},
			{0, 4, 5, 1},
			{1, 5, 6, 2},
			{2, 6, 7, 3},
			{3, 7, 4, 0},
			{4, 7, 6, 5},
		},
	}

	// Oktaeder - besteht aus 8 Dreiecken.
	Oktaeder = ObjectData{
		Points: []Vector{
			Vector{0, 0, -1},
			Vector{0, -1, 0},
			Vector{+1, 0, 0},
			Vector{0, +1, 0},
			Vector{-1, 0, 0},
			Vector{0, 0, +1},
		},
		SimpleFaces: [][]int{
			{0, 2, 1},
			{0, 3, 2},
			{0, 4, 3},
			{0, 1, 4},
			{1, 2, 5},
			{2, 3, 5},
			{3, 4, 5},
			{4, 1, 5},
		},
	}

	// Dodekaeder - besteht aus 12 regelmaessigen Fuenfecken.
	Dodekaeder = ObjectData{
		Points: []Vector{
			Vector{+phiI, 0, -phi}, // 0
			Vector{-phiI, 0, -phi}, // 1
			Vector{+1, +1, -1},     // 2
			Vector{-1, +1, -1},     // 3
			Vector{-1, -1, -1},     // 4
			Vector{+1, -1, -1},     // 5
			Vector{0, +phi, -phiI}, // 6
			Vector{0, -phi, -phiI}, // 7
			Vector{+phi, -phiI, 0}, // 8
			Vector{+phi, +phiI, 0}, // 9
			Vector{-phi, +phiI, 0}, // 10
			Vector{-phi, -phiI, 0}, // 11
			Vector{0, +phi, +phiI}, // 12
			Vector{0, -phi, +phiI}, // 13
			Vector{+1, +1, +1},     // 14
			Vector{-1, +1, +1},     // 15
			Vector{-1, -1, +1},     // 16
			Vector{+1, -1, +1},     // 17
			Vector{+phiI, 0, +phi}, // 18
			Vector{-phiI, 0, +phi}, // 19
		},
		SimpleFaces: [][]int{
			{0, 2, 9, 8, 5},
			{1, 4, 11, 10, 3},
			{0, 1, 3, 6, 2},
			{0, 5, 7, 4, 1},
			{2, 6, 12, 14, 9},
			{3, 10, 15, 12, 6},
			{5, 8, 17, 13, 7},
			{4, 7, 13, 16, 11},
			{8, 9, 14, 18, 17},
			{10, 11, 16, 19, 15},
			{12, 15, 19, 18, 14},
			{13, 17, 18, 19, 16},
		},
	}

	// Ikosaeder - besteht aus 20 (!) Dreiecken.
	Ikosaeder = ObjectData{
		Points: []Vector{
			Vector{0, -1, -phi}, // 0
			Vector{0, +1, -phi}, // +1
			Vector{-phi, 0, -1}, // 2
			Vector{+phi, 0, -1}, // 3
			Vector{-1, -phi, 0}, // 4
			Vector{+1, -phi, 0}, // 5
			Vector{+1, +phi, 0}, // 6
			Vector{-1, +phi, 0}, // 7
			Vector{-phi, 0, +1}, // 8
			Vector{+phi, 0, +1}, // 9
			Vector{0, -1, +phi}, // +10
			Vector{0, +1, +phi}, // +11
		},
		SimpleFaces: [][]int{
			{0, 1, 3},
			{0, 3, 5},
			{0, 5, 4},
			{0, 4, 2},
			{0, 2, 1},
			{1, 2, 7},
			{2, 8, 7},
			{2, 4, 8},
			{4, 10, 8},
			{4, 5, 10},
			{5, 9, 10},
			{5, 3, 9},
			{3, 6, 9},
			{3, 1, 6},
			{1, 7, 6},
			{8, 10, 11},
			{10, 9, 11},
			{9, 6, 11},
			{6, 7, 11},
			{7, 8, 11},
		},
	}

	//--------------------------------------------------------------

	// HollowCube -- Wuerfel mit ganz viel Luft darin
	b = 0.6

	HollowCube = ObjectData{
		Points: []Vector{
			// Die Eckpunkte des Wuerfels
			Vector{-1, -1, -1}, // 0
			Vector{-1, +1, -1},
			Vector{+1, +1, -1},
			Vector{+1, -1, -1},
			Vector{-1, -1, +1}, // 4
			Vector{+1, -1, +1},
			Vector{+1, +1, +1},
			Vector{-1, +1, +1},
			// Loch in der Bodenflaeche
			Vector{-b, -b, -1}, // 8
			Vector{-b, +b, -1},
			Vector{+b, +b, -1},
			Vector{+b, -b, -1},
			// Loch in der Deckflaeche
			Vector{-b, -b, +1}, // 12
			Vector{+b, -b, +1},
			Vector{+b, +b, +1},
			Vector{-b, +b, +1},
			// Loch in der linken Seitenflaeche
			Vector{-b, -1, -b}, // 16
			Vector{+b, -1, -b},
			Vector{+b, -1, +b},
			Vector{-b, -1, +b},
			// Loch in der rechten Seitenflaeche
			Vector{+b, +1, -b}, // 20
			Vector{-b, +1, -b},
			Vector{-b, +1, +b},
			Vector{+b, +1, +b},
			// Loch in der vorderen Seitenflaeche
			Vector{+1, -b, -b}, // 24
			Vector{+1, +b, -b},
			Vector{+1, +b, +b},
			Vector{+1, -b, +b},
			// Loch in der hinteren Seitenflaeche
			Vector{-1, +b, -b}, // 28
			Vector{-1, -b, -b},
			Vector{-1, -b, +b},
			Vector{-1, +b, +b},
			// Die Eckpunkte des inneren Wuerfels
			Vector{-b, -b, -b}, // 32
			Vector{-b, +b, -b},
			Vector{+b, +b, -b},
			Vector{+b, -b, -b},
			Vector{-b, -b, +b}, // 36
			Vector{+b, -b, +b},
			Vector{+b, +b, +b},
			Vector{-b, +b, +b},
		},
		TrigFaces: []FaceData{
			// Der Rand um das obere Loch
			{
				Edges: [][]int{
					{0, 3, 2, 1, 0}, {8, 11, 10, 9, 8},
				},
				Trigs: []TrigData{
					{0, 8, 3}, {3, 8, 11}, {3, 11, 2}, {2, 11, 10},
					{2, 10, 1}, {1, 10, 9}, {0, 1, 9}, {0, 9, 8},
				},
			},
			// Der Rand um das untere Loch
			{
				Edges: [][]int{
					{4, 5, 6, 7, 4}, {12, 13, 14, 15, 12},
				},
				Trigs: []TrigData{
					{4, 5, 12}, {5, 13, 12}, {5, 6, 13}, {6, 14, 13},
					{6, 7, 14}, {7, 15, 14}, {4, 15, 7}, {4, 12, 15},
				},
			},
			// Der Rand um das linke Loch
			{
				Edges: [][]int{
					{0, 3, 5, 4, 0}, {16, 17, 18, 19, 16},
				},
				Trigs: []TrigData{
					{4, 0, 16}, {4, 16, 19}, {4, 19, 18}, {4, 18, 5},
					{3, 16, 0}, {3, 17, 16}, {3, 18, 17}, {3, 5, 18},
				},
			},
			// Der Rand um das rechte Loch
			{
				Edges: [][]int{
					{1, 7, 6, 2, 1}, {20, 21, 22, 23, 20},
				},
				Trigs: []TrigData{
					{6, 2, 20}, {6, 20, 23}, {6, 23, 22}, {6, 22, 7},
					{1, 20, 2}, {1, 21, 20}, {1, 22, 21}, {1, 7, 22},
				},
			},
			// Der Rand um das vordere Loch
			{
				Edges: [][]int{
					{2, 6, 5, 3, 2}, {24, 25, 26, 27, 24},
				},
				Trigs: []TrigData{
					{5, 3, 24}, {5, 24, 27}, {5, 27, 26}, {5, 26, 6},
					{2, 24, 3}, {2, 25, 24}, {2, 26, 25}, {2, 6, 26},
				},
			},
			// Der Rand um das hintere Loch
			{
				Edges: [][]int{
					{0, 4, 7, 1, 0}, {28, 29, 30, 31, 28},
				},
				Trigs: []TrigData{
					{7, 1, 28}, {7, 28, 31}, {7, 31, 30}, {7, 30, 4},
					{0, 28, 1}, {0, 29, 28}, {0, 30, 29}, {0, 4, 30},
				},
			},
		},
		SimpleFaces: [][]int{
			// Die 4 Innenflaechen des...
			// ... unteren Loches
			{8, 32, 35, 11},
			{11, 35, 34, 10},
			{10, 34, 33, 9},
			{9, 33, 32, 8},
			// ... oberen Loches
			{12, 13, 37, 36},
			{13, 14, 38, 37},
			{14, 15, 39, 38},
			{15, 12, 36, 39},
			// ... linken Loches
			{16, 17, 35, 32},
			{17, 18, 37, 35},
			{18, 19, 36, 37},
			{19, 16, 32, 36},
			// ... rechten Loches
			{20, 21, 33, 34},
			{21, 22, 39, 33},
			{22, 23, 38, 39},
			{23, 20, 34, 38},
			// ... vorderen Loches
			{24, 25, 34, 35},
			{25, 26, 38, 34},
			{26, 27, 37, 38},
			{27, 24, 35, 37},
			// ... hinteren Loches
			{28, 29, 32, 33},
			{29, 30, 36, 32},
			{30, 31, 39, 36},
			{31, 28, 33, 39},
		},
	}

	// SwissFlag -- eine Art Schweizerflagge, wobei das Kreuz aus dem
	// Koerper entfernt wurde.
	c = 0.77
	d = 0.23
	e = 0.23

	SwissFlag = ObjectData{
		Points: []Vector{
			Vector{+1, +1, -e},
			Vector{-1, +1, -e},
			Vector{-1, -1, -e},
			Vector{+1, -1, -e},

			Vector{+1, +1, +e},
			Vector{-1, +1, +e},
			Vector{-1, -1, +e},
			Vector{+1, -1, +e},

			Vector{-c, -d, -e},
			Vector{-c, +d, -e},
			Vector{-d, +d, -e},
			Vector{-d, +c, -e},
			Vector{+d, +c, -e},
			Vector{+d, +d, -e},
			Vector{+c, +d, -e},
			Vector{+c, -d, -e},
			Vector{+d, -d, -e},
			Vector{+d, -c, -e},
			Vector{-d, -c, -e},
			Vector{-d, -d, -e},

			Vector{-c, -d, +e},
			Vector{-c, +d, +e},
			Vector{-d, +d, +e},
			Vector{-d, +c, +e},
			Vector{+d, +c, +e},
			Vector{+d, +d, +e},
			Vector{+c, +d, +e},
			Vector{+c, -d, +e},
			Vector{+d, -d, +e},
			Vector{+d, -c, +e},
			Vector{-d, -c, +e},
			Vector{-d, -d, +e},
		},
		TrigFaces: []FaceData{
			// Die Deckflaeche, welche das Schweizerkreuz als Aussparung
			// hat.
			{
				Edges: [][]int{
					{0, 1, 2, 3, 0},
					{8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 8},
				},
				Trigs: []TrigData{
					{0, 13, 12}, {0, 14, 13},
					{0, 12, 1}, {1, 12, 11},
					{1, 10, 9}, {1, 11, 10},
					{1, 9, 2}, {2, 9, 8},
					{2, 19, 18}, {2, 8, 19},
					{2, 18, 3}, {3, 18, 17},
					{3, 16, 15}, {3, 17, 16},
					{3, 15, 0}, {0, 15, 14},
				},
			},
			// Die Bodenflaeche, welche das Schweizerkreuz als Aussparung
			// hat.
			{
				Edges: [][]int{
					{4, 5, 6, 7, 4},
					{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 20},
				},
				Trigs: []TrigData{
					{4, 24, 25}, {4, 25, 26},
					{4, 5, 24}, {5, 23, 24},
					{5, 21, 22}, {5, 22, 23},
					{5, 6, 21}, {6, 20, 21},
					{6, 30, 31}, {6, 31, 20},
					{6, 7, 30}, {7, 29, 30},
					{7, 27, 28}, {7, 28, 29},
					{7, 4, 27}, {4, 26, 27},
				},
			},
		},
		SimpleFaces: [][]int{
			// Die rundumlaufenden aeusseren Begrenzungsflaechen
			{0, 1, 5, 4},
			{1, 2, 6, 5},
			{2, 3, 7, 6},
			{3, 0, 4, 7},
			// Die inneren Begrenzungsflaechen entlang des Kreuzes
			{8, 9, 21, 20},
			{9, 10, 22, 21},
			{10, 11, 23, 22},
			{11, 12, 24, 23},
			{12, 13, 25, 24},
			{13, 14, 26, 25},
			{14, 15, 27, 26},
			{15, 16, 28, 27},
			{16, 17, 29, 28},
			{17, 18, 30, 29},
			{18, 19, 31, 30},
			{19, 8, 20, 31},
		},
	}

	// 3-dimensionales (Schweizer)kreuz
	f = 0.3

	SwissCross = ObjectData{
		Points: []Vector{
			Vector{-f, -f, -1}, // 0
			Vector{+f, -f, -1},
			Vector{+f, +f, -1},
			Vector{-f, +f, -1},

			Vector{-f, -f, -f}, // 4
			Vector{+f, -f, -f},
			Vector{+f, +f, -f},
			Vector{-f, +f, -f},

			Vector{-f, -1, -f}, // 8
			Vector{+f, -1, -f},
			Vector{+1, -f, -f},
			Vector{+1, +f, -f},
			Vector{+f, +1, -f},
			Vector{-f, +1, -f},
			Vector{-1, +f, -f},
			Vector{-1, -f, -f},

			Vector{-f, -1, +f}, // 16
			Vector{+f, -1, +f},
			Vector{+1, -f, +f},
			Vector{+1, +f, +f},
			Vector{+f, +1, +f},
			Vector{-f, +1, +f},
			Vector{-1, +f, +f},
			Vector{-1, -f, +f},

			Vector{-f, -f, +f}, // 24
			Vector{+f, -f, +f},
			Vector{+f, +f, +f},
			Vector{-f, +f, +f},

			Vector{-f, -f, +1}, // 28
			Vector{+f, -f, +1},
			Vector{+f, +f, +1},
			Vector{-f, +f, +1},
		},
		SimpleFaces: [][]int{
			{0, 3, 2, 1},
			{0, 1, 5, 4}, {1, 2, 6, 5}, {2, 3, 7, 6}, {3, 0, 4, 7},
			{4, 5, 9, 8}, {5, 6, 11, 10}, {6, 7, 13, 12}, {7, 4, 15, 14},
			{8, 9, 17, 16}, {10, 11, 19, 18}, {12, 13, 21, 20},
			{14, 15, 23, 22},
			{9, 5, 25, 17}, {5, 10, 18, 25}, {11, 6, 26, 19}, {6, 12, 20, 26},
			{13, 7, 27, 21}, {7, 14, 22, 27}, {15, 4, 24, 23}, {4, 8, 16, 24},
			{25, 24, 16, 17}, {26, 25, 18, 19}, {27, 26, 20, 21},
			{24, 27, 22, 23},
			{24, 25, 29, 28}, {25, 26, 30, 29}, {26, 27, 31, 30},
			{27, 24, 28, 31},
			{28, 29, 30, 31},
		},
	}

	// Spezial-Objekt, um die Triangulation von nicht-konvexen Flaechen
	// zu testen.

	g = 0.3
	h = 0.4
	i = 0.5

	TestBody = ObjectData{
		Points: []Vector{
			Vector{-1, -1, -1}, // 0
			Vector{+1, -1, -1},
			Vector{+1, +1, -1},
			Vector{-1, +1, -1},

			Vector{-1, -h, -g}, // 4
			Vector{+1, -h, -g},
			Vector{+1, +1, -i},
			Vector{-1, +1, -i},

			Vector{-1, -h, +g}, // 8
			Vector{-g, -h, +g},
			Vector{+g, -h, +g},
			Vector{+1, -h, +g},
			Vector{+1, +1, +i},
			Vector{+i, +1, +i},
			Vector{-i, +1, +i},
			Vector{-1, +1, +i},

			Vector{-1, -1, +1}, // 16
			Vector{+1, -1, +1},
			Vector{+1, +1, +1},
			Vector{+i, +1, +1},
			Vector{+g, -h, +1},
			Vector{-g, -h, +1},
			Vector{-i, +1, +1},
			Vector{-1, +1, +1},

			Vector{-1, -h, +1}, // 24
		},
		SimpleFaces: [][]int{
			{0, 3, 2, 1},
			{2, 3, 7, 6},
			{6, 7, 4, 5},
			//{4, 24, 20, 10, 11, 5},
			{11, 10, 13, 12},
			//{9, 8, 15, 14},
			{12, 13, 19, 18},
			{13, 10, 20, 19},
			//{9, 14, 22, 21},
			//{14, 15, 23, 22},
			{17, 18, 19, 20, 24, 16},
			//{17, 18, 19, 20, 21, 22, 23, 16},
			{0, 1, 17, 16},
			// hier muesste es jetzt knallen...
			{1, 2, 6, 5, 11, 12, 18, 17},
			{16, 24, 4, 7, 3, 0},
			//{16, 23, 15, 8, 4, 7, 3, 0},
		},
		TrigFaces: []FaceData{
			{
				Edges: [][]int{
					{4, 24, 20, 10, 11, 5},
				},
				Trigs: []TrigData{
					{5, 10, 11}, {5, 4, 10}, {10, 4, 8},
					{24, 20, 10}, {10, 8, 24},
				},
			},
		},
	}

	HiddenFaces = ObjectData{
		Points: []Vector{
			Vector{-1, -1, -1}, // 0
			Vector{+1, -1, -1},
			Vector{+1, +1, -1},
			Vector{-1, +1, -1},

			Vector{-1, -h, -g}, // 4
			Vector{+1, -h, -g},
			Vector{+1, +1, -i},
			Vector{-1, +1, -i},

			Vector{-1, -h, +g}, // 8
			Vector{-g, -h, +g},
			Vector{+g, -h, +g},
			Vector{+1, -h, +g},
			Vector{+1, +1, +i},
			Vector{+i, +1, +i},
			Vector{-i, +1, +i},
			Vector{-1, +1, +i},

			Vector{-1, -1, +1}, // 16
			Vector{+1, -1, +1},
			Vector{+1, +1, +1},
			Vector{+i, +1, +1},
			Vector{+g, -h, +1},
			Vector{-g, -h, +1},
			Vector{-i, +1, +1},
			Vector{-1, +1, +1},

			Vector{-1, -h, +1}, // 24

			Vector{+g, -h, -g}, // 25 Hilfspunkt fuer Triangulation
		},
		SimpleFaces: [][]int{
			{6, 7, 4, 5},
			//{4, 24, 20, 10, 11, 5},
			{11, 10, 13, 12},
			{12, 13, 19, 18},
			{13, 10, 20, 19},
		},
		TrigFaces: []FaceData{
			{
				Edges: [][]int{
					{4, 24, 20, 10, 11, 5},
				},
				Trigs: []TrigData{
					{5, 25, 10}, {5, 10, 11}, {4, 24, 20}, {4, 20, 25},
				},
			},
		},
	}

	j = 0.5

	HiddenFaces2 = ObjectData{
		Points: []Vector{
			Vector{-1, -1, -1},
			Vector{0, -1 / 2, -1},
			Vector{+1, -1, -1},
			Vector{+1 / 2, 0, -1},
			Vector{+1, +1, -1},
			Vector{0, +1 / 2, -1},
			Vector{-1, +1, -1},
			Vector{-1 / 2, 0, -1},
			Vector{-j, -j, 0},
			Vector{+j, -j, 0},
			Vector{+j, +j, 0},
			Vector{-j, +j, 0},
			Vector{-0.8, -0.8, 1},
			Vector{-0.3, -0.6, 1},
			Vector{+0.2, -0.8, 1},
			Vector{+0.8, -0.6, 1},
			Vector{+0.8, -0.1, 1},
			Vector{+0.5, -0.2, 1},
			Vector{+0.2, 0.0, 1},
			Vector{+0.3, +0.3, 1},
			Vector{+0.8, +0.3, 1},
			Vector{+0.8, +0.8, 1},
			Vector{+0.2, +0.8, 1},
			Vector{-0.2, +0.5, 1},
			Vector{-0.8, +0.8, 1},
			Vector{-0.6, +0.1, 1},
			Vector{0.0, 0.0, 1},
			Vector{-0.6, -0.1, 1},
		},
		SimpleFaces: [][]int{
			{0, 1, 2, 3, 4, 5, 6, 7},
			{8, 9, 10, 11},
			//{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
		},
	}
)
