package assets

import (
	"bytes"
	"embed"
	_ "embed"
	"image"
	_ "image/png"
)

//go:embed appicon.png

// AppIcon is encoded protonet's icon in png format
var AppIcon []byte

// AppIconImage is decoded Image representing AppIcon
var AppIconImage, _, _ = image.Decode(bytes.NewReader(AppIcon))

//go:embed appicon.png
var AppIconFile embed.FS
var AppIconFileName = "appicon.png"
