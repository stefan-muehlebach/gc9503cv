//go:build ignore

//----------------------------------------------------------------------------

type Canvas struct {
	Objs *list.List
	Anims *list.List
	Rect image.Rectangle
	GC *gg.Context
	Img draw.Image
	BackColor colors.RGBA
}

// Nicht exportierte Funktion! Nur der Screen (GC9503CV) soll solche Objekte
// erstellen duerfen.
func newCanvas(size geom.Point[int]) *Canvas {
	c := &Canvas{}
	c.Objs = list.New()
	c.Anims = list.New()
	c.Rect = image.Rect(0, 0, size.X, size.Y)
	c.BackColor = colors.Transparent
	c.Img = image.NewRGBA(c.Rect)
	c.GC = gg.NewContextForRGBA(c.Img.(*image.RGBA))
	return c
}

func (c *Canvas) Add(objs ...Node) {
    for _, obj := range objs {
		c.Objs.PushBack(obj)
	}
}

func (c *Canvas) Del(obj Node) {
	for ele := c.Objs.Front(); ele != nil; ele = ele.Next() {
		o := ele.Value.(Node)
		if o == obj {
			c.Objs.Remove(ele)
			return
		}
	}
}

func (c *Canvas) Purge() {
	c.Objs.Init()
}

func (c *Canvas) FindTarget(pt Point) Node {
	for ele := c.Objs.Front(); ele != nil; ele = ele.Next() {
		obj := ele.Value.(Node)
		if !obj.IsVisible() {
			continue
		}
		if target := obj.Contains(pt); target != nil {
			return target
		}
	}
	return nil
}

func (c *Canvas) Clear(color colors.RGBA) {
	draw.Draw(c.Img, c.Rect, image.NewUniform(color), image.Point{}, draw.Src)
}

func (c *Canvas) Refresh() {
	var obj Node
	var ok bool

	c.Clear(c.BackColor)
	for ele := c.Objs.Front(); ele != nil; ele = ele.Next() {
		if obj, ok = ele.Value.(Node); !ok {
			log.Printf("wrong object in object list of canvas")
			continue
		}
		if !obj.IsVisible() {
			continue
		}
		obj.Draw(c.GC)
	}
}

