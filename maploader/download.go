package maploader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func download(workerIndex int, url URL) error {
	fmt.Printf(" * [%03d] Downloading %s\n", workerIndex, url.In)
	err := os.MkdirAll(filepath.Dir(url.Out), os.ModePerm)
	if err != nil {
		return err
	}

	outfile, err := os.Create(url.Out)
	if err != nil {
		return err
	}
	defer outfile.Close()

	resp, err := http.Get(url.In)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Worker downloads maps
func Worker(index int, urlChannel chan URL) {
	fmt.Printf("Starting worker %03d\n", index)
	for {
		url, ok := <-urlChannel
		if ok {
			err := download(index, url)
			if err != nil {
				fmt.Printf("FAIL to download %s\n", url.In)
			}
		} else {
			break
		}
	}
}
