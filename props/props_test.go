package props

import (
	"testing"

	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
)

var (
	defProps      *Properties
	typeProps     *Properties
	objProps      *Properties
	c1, c2, c3    colors.RGBA
	f             *fonts.Font
	s1            float64
	colorPropName = Color
)

func init() {
	defProps = PropsMap["Default"]
	typeProps = PropsMap["Button"]

	objProps = NewProperties(typeProps)

}

func TestPropertyTree(t *testing.T) {
	t.Logf("Object color: %+v", objProps.Color(colorPropName))
}

// Prüft, ob in den Default-Properties zu allen Property-Namen ein Eintrag
// vorhanden ist.
func TestDefaultProperties(t *testing.T) {
	for name := ColorPropertyName(0); name < NumColorProperties; name++ {
		_, ok := defProps.ColorMap[name]
		if !ok {
			t.Errorf("Color property '%+v' not in defaults", name)
		}
	}

	for name := FontPropertyName(0); name < NumFontProperties; name++ {
		_, ok := defProps.FontMap[name]
		if !ok {
			t.Errorf("Font property '%+v' not in defaults", name)
		}
	}

	for name := SizePropertyName(0); name < NumSizeProperties; name++ {
		_, ok := defProps.SizeMap[name]
		if !ok {
			t.Errorf("Size property '%+v' not in defaults", name)
		}
	}
}

// Prüft, ob die hierarchische Vererbung über drei Stufen (Default, Type,
// Object) funktioniert, übersteuert und wieder gelöscht werden kann.
func TestColorHierarchy(t *testing.T) {
	c1 = defProps.Color(colorPropName)
	c2 = typeProps.Color(colorPropName)
	c3 = objProps.Color(colorPropName)
	t.Logf("Default color   : %+v", c1)
	t.Logf("  Type color    : %+v", c2)
	t.Logf("    Object color: %+v", c3)

	if c2 != c1 {
		t.Errorf("Default and type prop differ (got '%v', want '%v'", c2, c1)
	}
	if c3 != c1 {
		t.Errorf("Default and object prop differ (got '%v', want '%v'", c3, c1)
	}

	typeProps.SetColor(colorPropName, colors.FireBrick)
	objProps.SetColor(colorPropName, colors.Yellow)

	c1 = defProps.Color(colorPropName)
	c2 = typeProps.Color(colorPropName)
	c3 = objProps.Color(colorPropName)
	t.Logf("Default color   : %+v", c1)
	t.Logf("  Type color    : %+v", c2)
	t.Logf("    Object color: %+v", c3)

	if c2 == c1 {
		t.Errorf("Default and type prop are equal (got '%v', want '%v'", c2, c1)
	}
	if c3 == c1 {
		t.Errorf("Default and object prop are equal (got '%v', want '%v'", c3, c1)
	}

	typeProps.DelColor(colorPropName)
	objProps.DelColor(colorPropName)

	c1 = defProps.Color(colorPropName)
	c2 = typeProps.Color(colorPropName)
	c3 = objProps.Color(colorPropName)
	t.Logf("Default color   : %+v", c1)
	t.Logf("  Type color    : %+v", c2)
	t.Logf("    Object color: %+v", c3)

	if c2 != c1 {
		t.Errorf("Default and type prop differ (got '%v', want '%v'", c2, c1)
	}
	if c3 != c1 {
		t.Errorf("Default and object prop differ (got '%v', want '%v'", c3, c1)
	}
}

func TestGetFont(t *testing.T) {
	f = defProps.Font(BoldFont)
	t.Logf("Def.BoldFont: %T", f)
}

// Direkter Zugriff auf Default-Property.
func BenchmarkGetDefColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c1 = defProps.Color(BorderColor)
	}
}

// Indirekter Zugriff via Type-Property.
func BenchmarkGetTypeColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c1 = typeProps.Color(BorderColor)
	}
}

// Doppelt indirekter Zugriff via Object-Property.
func BenchmarkGetObjColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c1 = objProps.Color(BorderColor)
	}
}

// Direkter Zugriff auf Default-Property.
func BenchmarkGetDefSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 = defProps.Size(BorderWidth)
	}
}

// Indirekter Zugriff via Type-Property.
func BenchmarkGetTypeSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 = typeProps.Size(BorderWidth)
	}
}

// Doppelt indirekter Zugriff via Object-Property.
func BenchmarkGetObjSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1 = objProps.Size(BorderWidth)
	}
}

func TestPropFileRead(t *testing.T) {
	PropsMap = NewPropsMapFromUserFile("TestProps.json")
	t.Logf("Default.Color: %+v", PropsMap["Default"].Color(Color))
}
