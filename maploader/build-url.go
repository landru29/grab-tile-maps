package maploader

import (
	"fmt"
	"math"
	"regexp"
)

func queueUrlsByZoom(inFolderScheme string, outFolderScheme string, zoom int, queue chan URL) {
	rx := regexp.MustCompile("{x}")
	ry := regexp.MustCompile("{y}")
	rz := regexp.MustCompile("{z}")
	size := int64(math.Pow(2, float64(zoom)))
	fmt.Printf("  * zoom %d ", zoom)
	for y := int64(0); y < size; y++ {
		for x := int64(0); x < size; x++ {
			url := rx.ReplaceAllString(inFolderScheme, fmt.Sprintf("%d", x))
			url = ry.ReplaceAllString(url, fmt.Sprintf("%d", y))
			url = rz.ReplaceAllString(url, fmt.Sprintf("%d", zoom))
			out := rx.ReplaceAllString(
				outFolderScheme,
				fmt.Sprintf("%d", x),
			)
			out = ry.ReplaceAllString(out, fmt.Sprintf("%d", y))
			out = rz.ReplaceAllString(out, fmt.Sprintf("%d", zoom))
			queue <- URL{
				In:  url,
				Out: out,
			}
		}
	}
}

// QueueUrls pushs urls in the download queue
func QueueUrls(inFolderScheme string, outFolderScheme string, zoomMin int, zoomMax int, queue chan URL) {
	for zoom := zoomMin; zoom <= zoomMax; zoom++ {
		queueUrlsByZoom(inFolderScheme, outFolderScheme, zoom, queue)
	}
}
