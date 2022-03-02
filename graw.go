package graw

import (
	"encoding/xml"
	"strings"
)

const (
	topCellId  = "0"
	rootCellID = "1"
)

// GraphModel 默认有一个顶级cell，其他cell parentId指向1
//     <mxCell id="0" style=";html=1;" />
//     <mxCell id="1" style=";html=1;" parent="0" />
type GraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Dx      int      `xml:"dx,attr"`
	Dy      int      `xml:"dy,attr"`

	// 属性
	Grid       string `xml:"grid,attr,omitempty"`
	GridSize   string `xml:"gridSize,attr,omitempty"`
	Guides     string `xml:"guides,attr,omitempty"`
	Tooltips   string `xml:"tooltips,attr,omitempty"`
	Connect    string `xml:"connect,attr,omitempty"`
	Arrows     string `xml:"arrows,attr,omitempty"`
	Fold       string `xml:"fold,attr,omitempty"`
	Page       string `xml:"page,attr,omitempty"`
	PageScale  string `xml:"pageScale,attr,omitempty"`
	PageWidth  string `xml:"pageWidth,attr,omitempty"`
	PageHeight string `xml:"pageHeight,attr,omitempty"`
	Background string `xml:"background,attr,omitempty"`
	Math       string `xml:"math,attr,omitempty"`
	Shadow     string `xml:"shadow,attr,omitempty"`

	Root []Cell `xml:"root>mxCell"`
}

// Cell 单元格/元素
// Vertex=1 为顶点
// Edge=1 为边
// ParentID 为根单元格的ID
type Cell struct {
	XMLName  xml.Name `xml:"mxCell"`
	ID       string   `xml:"id,attr"`
	Value    string   `xml:"value,attr,omitempty"`
	Style    Style    `xml:"style,attr,omitempty"`
	ParentID string   `xml:"parent,attr,omitempty"`
	Vertex   string   `xml:"vertex,attr,omitempty"`
	Edge     string   `xml:"edge,attr,omitempty"`
	Source   string   `xml:"source,attr,omitempty"`
	Target   string   `xml:"target,attr,omitempty"`
	Geometry *Geometry
}

// Geometry
// 意为测量土地的学问。这里我理解为对图形位置的描述
// 初始坐标 X Y
// 宽高 Width Height
//
// as 默认为 "geometry"
type Geometry struct {
	XMLName  xml.Name `xml:"mxGeometry"`
	X        int      `xml:"x,attr,omitempty"`
	Y        int      `xml:"y,attr,omitempty"`
	Width    string   `xml:"width,attr,omitempty"`
	Height   string   `xml:"height,attr,omitempty"`
	Relative string   `xml:"relative,attr,omitempty"`
	As       string   `xml:"as,attr"`
	Point    *Point
}

// Point
// 几何图形边缘的点坐标
// as 的值 ["sourcePoint","targetPoint"]
// 对线来说就是起点和终点

type Point struct {
	XMLName xml.Name `xml:"mxPoint"`
	X       int      `xml:"x,attr,omitempty"`
	Y       int      `xml:"y,attr,omitempty"`
	As      string   `xml:"as,attr"`
}

// A Style is a map of key-value pairs to describe the style
// properties of each cell.
type Style struct {
	Attributes map[string]string
}

// MarshalXMLAttr returns an XML attribute with the encoded value
// of Style. It implements xml.MarshalerAttr interface.
func (a Style) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var text string

	for k, v := range a.Attributes {
		text += k

		if v != "" {
			text += "="
			text += v
		}

		text += ";"
	}

	return xml.Attr{Name: xml.Name{Local: "style"}, Value: text}, nil
}

// UnmarshalXMLAttr decodes a single XML attribute of type Style.
// It implements xml.UnmarshalerAttr interface.
func (a *Style) UnmarshalXMLAttr(attr xml.Attr) error {
	a.Attributes = make(map[string]string)
	pairs := strings.Split(attr.Value, ";")

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) < 2 {
			kv = append(kv, "")
		}
		a.Attributes[kv[0]] = kv[1]
	}

	return nil
}

// NewGraph returns a new graph model containing a root cell and
// one layer with ID layerId.
func NewGraph() GraphModel {
	rootStyle := Style{
		Attributes: make(map[string]string),
	}
	rootStyle.Attributes["html"] = "1"
	return GraphModel{
		Dx: 640,
		Dy: 480,
		Root: []Cell{
			{ID: topCellId, Style: rootStyle},
			{
				ID:       rootCellID,
				ParentID: topCellId,
				Style:    rootStyle,
			},
		},
	}
}

// Add adds the given Cell to the root cell of the receiving
// graph model.
func (g *GraphModel) Add(c *Cell) *GraphModel {
	g.Root = append(g.Root, *c)
	return g
}

// NewShape returns a new Vertex Cell, configured with the given
// unique ID (id) and parent ID (layerId). The new cell contains
// a default geometry which you might want to change.
func NewShape(id, layerId string) *Cell {
	s := newCell(id, layerId)
	s.Vertex = "1"
	s.Geometry = newGeometry()
	return s
}

// NewImage returns a new Vertex Cell, configured as an image,
// with given unique ID (id), parent ID (layerId) and image
// source URL (url). The new cell contains a default geometry
// which you might want to change.
func NewImage(id, layerId, url string) *Cell {
	i := NewShape(id, layerId)
	i.Style = Style{
		Attributes: map[string]string{
			"shape":       "image",
			"imageAspect": "0",
			"image":       url,
		},
	}
	return i
}

// NewImageXY returns a new Vertex Cell, configured as an image,
// with given unique ID (id), parent ID (layerId), image source
// URL (url), and coordinates (x and y).
func NewImageXY(id, layerId, url string, x int, y int) *Cell {
	i := NewImage(id, layerId, url)
	i.Geometry.X = x
	i.Geometry.X = y
	return i
}

// newCell returns a new Cell object configured with id and
// parent ID.
func newCell(id string, layerId string) *Cell {
	return &Cell{
		ID:       id,
		ParentID: layerId,
	}

}

// newGeometry returns a new Geometry object configured with
// default values.
func newGeometry() *Geometry {
	return &Geometry{
		X:  10,
		Y:  10,
		As: "geometry",
	}
}
