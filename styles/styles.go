package styles

import "github.com/tealeg/xlsx"

func Wrap(c *xlsx.Cell) *xlsx.Cell {
	s := c.GetStyle()
	a := &s.Alignment
	a.WrapText = true
	s.Alignment = *a
	s.ApplyAlignment = true
	c.SetStyle(s)
	return c
}

func Center(c *xlsx.Cell) *xlsx.Cell {
	s := c.GetStyle()
	a := &s.Alignment
	a.Horizontal = "center"
	a.Vertical = "center"
	s.Alignment = *a
	s.ApplyAlignment = true
	c.SetStyle(s)
	return c
}

func Header(c *xlsx.Cell) *xlsx.Cell {
	c = Center(Wrap(c))
	b := xlsx.NewBorder("thin", "thin", "thin", "thin")
	b.TopColor = "bbbbbb"
	b.LeftColor = "bbbbbb"
	b.RightColor = "bbbbbb"
	b.BottomColor = "bbbbbb"
	font := xlsx.NewFont(12, "arial")
	font.Bold = true
	s := c.GetStyle()
	s.Fill = *xlsx.NewFill("solid", "e5e5e5", "333333")
	s.Font = *font
	s.Border = *b
	s.ApplyFont = true
	s.ApplyFill = true
	s.ApplyBorder = true
	return c
}

func Rate(c *xlsx.Cell) *xlsx.Cell {
	c = Center(Wrap(c))

	b := xlsx.NewBorder("thin", "thin", "thin", "thin")
	b.TopColor = "dddddd"
	b.LeftColor = "dddddd"
	b.RightColor = "dddddd"
	b.BottomColor = "dddddd"
	s := c.GetStyle()
	s.Fill = *xlsx.NewFill("solid", "D9FFEB", "00000000")
	s.Border = *b
	s.ApplyFont = true
	s.ApplyFill = true
	s.ApplyBorder = true
	return c
}

func Error(c *xlsx.Cell) *xlsx.Cell {
	s := c.GetStyle()
	s.Fill = *xlsx.NewFill("solid", "f2dede", "00000000")
	s.ApplyFont = true
	s.ApplyFill = true
	s.ApplyBorder = true
	return c
}
