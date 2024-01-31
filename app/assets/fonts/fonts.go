package fonts

import (
	_ "embed"
	"gioui.org/font"
	"github.com/partisiadev/partisiawallet/app/assets/fonts/roboto"
)

func unpackMultiFaces(faces [][]font.FontFace) []font.FontFace {
	newFaces := make([]font.FontFace, 0)
	for _, facesSlice := range faces {
		newFaces = append(newFaces, facesSlice...)
	}
	return newFaces
}

var Collection = append(
	unpackMultiFaces(roboto.MultiCollectionFaces),
)
