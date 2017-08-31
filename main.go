package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type URL struct {
	In  string
	Out string
}

func buildUrls(zoom int) []URL {
	rx := regexp.MustCompile("{x}")
	ry := regexp.MustCompile("{y}")
	rz := regexp.MustCompile("{z}")
	size := int64(math.Pow(2, float64(zoom)))
	var urls []URL
	for y := int64(0); y < size; y++ {
		for x := int64(0); x < size; x++ {
			url := rx.ReplaceAllString(viper.GetString("url_scheme"), fmt.Sprintf("%d", x))
			url = ry.ReplaceAllString(url, fmt.Sprintf("%d", y))
			url = rz.ReplaceAllString(url, fmt.Sprintf("%d", zoom))
			out := rx.ReplaceAllString(
				fmt.Sprintf("%s/%s", viper.GetString("out_folder"), viper.GetString("out_scheme")),
				fmt.Sprintf("%d", x),
			)
			out = ry.ReplaceAllString(out, fmt.Sprintf("%d", y))
			out = rz.ReplaceAllString(out, fmt.Sprintf("%d", zoom))
			urls = append(
				urls,
				URL{
					In:  url,
					Out: out,
				},
			)
		}
	}
	return urls
}

func download(url URL) error {
	fmt.Printf(" * Downloading %s ", url.In)
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

var mainCommand = &cobra.Command{
	Use:   "grab-tile-maps",
	Short: "grab-tile-maps",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.GetString("sql_user")

		for zoom := viper.GetInt("zoom_min"); zoom <= viper.GetInt("zoom_max"); zoom++ {
			fmt.Printf("Zoom %d:\n", zoom)
			urls := buildUrls(zoom)
			for _, url := range urls {
				err := download(url)
				if err != nil {
					fmt.Printf("FAIL\n")
				}
				fmt.Printf("OK\n")
			}
		}
	},
}

func init() {

	viper.SetEnvPrefix("gtm")
	viper.AutomaticEnv()
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	flags := mainCommand.Flags()

	flags.String("out-folder", "./map", "output folder")
	flags.Int("zoom-min", 0, "Minimum zoom")
	flags.Int("zoom-max", 3, "Maximum zoom")
	flags.String("url-scheme", "https://b.tile.opentopomap.org/{z}/{x}/{y}.png", "Map tile urls")
	flags.String("out-scheme", "{z}/{x}/{y}.png", "Map tile out files")

	viper.BindPFlag("out_folder", flags.Lookup("out-folder"))
	viper.BindPFlag("url_scheme", flags.Lookup("url-scheme"))
	viper.BindPFlag("zoom_min", flags.Lookup("zoom-min"))
	viper.BindPFlag("zoom_max", flags.Lookup("zoom-max"))
	viper.BindPFlag("out_scheme", flags.Lookup("out-scheme"))

}

func main() {
	mainCommand.Execute()
}
