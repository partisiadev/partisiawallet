package roboto

import (
	_ "embed"
	"gioui.org/font"
	"gioui.org/font/opentype"
)

//go:embed Roboto-Black.ttf
var BlackTTF []byte

//go:embed Roboto-BlackItalic.ttf
var BlackItalicTTF []byte

//go:embed Roboto-Bold.ttf
var BoldTTF []byte

//go:embed Roboto-BoldItalic.ttf
var BoldItalicTTF []byte

//go:embed Roboto-Italic.ttf
var ItalicTTF []byte

//go:embed Roboto-Light.ttf
var LightTTF []byte

//go:embed Roboto-LightItalic.ttf
var LightItalicTTF []byte

//go:embed Roboto-Medium.ttf
var MediumTTF []byte

//go:embed Roboto-MediumItalic.ttf
var MediumItalicTTF []byte

//go:embed Roboto-Regular.ttf
var RegularTTF []byte

//go:embed Roboto-Thin.ttf
var ThinTTF []byte

//go:embed Roboto-ThinItalic.ttf
var ThinItalicTTF []byte

var BlackFaces, _ = opentype.ParseCollection(BlackTTF)
var BlackItalicFaces, _ = opentype.ParseCollection(BlackItalicTTF)
var BoldFaces, _ = opentype.ParseCollection(BoldTTF)
var BoldItalicFaces, _ = opentype.ParseCollection(BoldItalicTTF)
var ItalicFaces, _ = opentype.ParseCollection(ItalicTTF)
var LightFaces, _ = opentype.ParseCollection(LightTTF)
var LightItalicFaces, _ = opentype.ParseCollection(LightItalicTTF)
var MediumFaces, _ = opentype.ParseCollection(MediumTTF)
var MediumItalicFaces, _ = opentype.ParseCollection(MediumItalicTTF)
var RegularFaces, _ = opentype.ParseCollection(RegularTTF)
var ThinFaces, _ = opentype.ParseCollection(ThinTTF)
var ThinItalicFaces, _ = opentype.ParseCollection(ThinItalicTTF)

var MultiCollectionFaces = [][]font.FontFace{
	BlackFaces,
	BlackItalicFaces,
	BoldFaces,
	BoldItalicFaces,
	ItalicFaces,
	LightFaces,
	LightItalicFaces,
	MediumFaces,
	MediumItalicFaces,
	RegularFaces,
	ThinFaces,
	ThinItalicFaces,
}
