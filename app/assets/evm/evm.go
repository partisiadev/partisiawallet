package evm

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/partisiadev/partisiawallet/log"
	"io/fs"
	"strings"
)

//go:embed chains
var ChainsDir embed.FS
var ChainsDirName = "chains"

//go:embed iconsDownload
var IconsDownloadDir embed.FS
var IconsDownloadDirName = "iconsDownload"

//go:embed icons
var IconsDir embed.FS
var IconsDirName = "icons"

var ErrUnsupportedImage = errors.New("unsupported image")

type LoadChainsResponse struct {
	Done   bool
	Chains []Chain
}

var loadChainsResponse = LoadChainsResponse{}

func LoadChains() <-chan LoadChainsResponse {
	if loadChainsResponse.Done {
		resCh := make(chan LoadChainsResponse, 1)
		resCh <- loadChainsResponse
		close(resCh)
		return resCh
	}
	files, err := ChainsDir.ReadDir(ChainsDirName)
	if err != nil {
		log.Logger().Fatal(err)
	}
	resCh := make(chan LoadChainsResponse, len(files))
	go func(files []fs.DirEntry, resCh chan<- LoadChainsResponse) {
		chains := make([]Chain, 0)
		var resp LoadChainsResponse
		resp.Chains = chains
		defer func() {
			resp.Done = true
			resp.Chains = chains
			loadChainsResponse = resp
			resCh <- resp
			close(resCh)
		}()
		for _, file := range files {
			var val []byte
			val, err = ChainsDir.ReadFile(ChainsDirName + "/" + file.Name())
			if err != nil {
				continue
			}
			var chain Chain
			err = json.Unmarshal(val, &chain)
			if err != nil {
				continue
			}
			// skip all deprecated chains
			if strings.ToLower(chain.Status) == "deprecated" {
				continue
			}
			iconsDataFile, err := IconsDir.ReadFile(IconsDirName + "/" + chain.Icon + ".json")
			if err != nil {
				continue
			}
			var iconsData []IconData
			err = json.Unmarshal(iconsDataFile, &iconsData)
			if err != nil {
				continue
			}
			chain.IconsData = iconsData
			chains = append(chains, chain)
			resp.Chains = chains
			resCh <- resp
		}
	}(files, resCh)
	return resCh
}
