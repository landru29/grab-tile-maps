package main

import (
	"fmt"

	"github.com/landru29/grab-tile-maps/maploader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCommand = &cobra.Command{
	Use:   "grab-tile-maps",
	Short: "grab-tile-maps",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		urlChannel := make(chan maploader.URL)

		for i := 0; i < viper.GetInt("download_agent"); i++ {
			go maploader.Worker(i, urlChannel)
		}

		maploader.QueueUrls(
			viper.GetString("url_scheme"),
			fmt.Sprintf("%s/%s", viper.GetString("out_folder"), viper.GetString("out_scheme")),
			viper.GetInt("zoom_min"),
			viper.GetInt("zoom_max"),
			urlChannel,
		)

		close(urlChannel)
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

	flags.String("out-folder", "./out", "output folder")
	flags.Int("zoom-min", 0, "Minimum zoom")
	flags.Int("zoom-max", 3, "Maximum zoom")
	flags.String("url-scheme", "https://b.tile.opentopomap.org/{z}/{x}/{y}.png", "Map tile urls")
	flags.String("out-scheme", "{z}/{x}/{y}.png", "Map tile out files")
	flags.Int("agent", 10, "Number of agents for downloading")

	viper.BindPFlag("out_folder", flags.Lookup("out-folder"))
	viper.BindPFlag("url_scheme", flags.Lookup("url-scheme"))
	viper.BindPFlag("zoom_min", flags.Lookup("zoom-min"))
	viper.BindPFlag("zoom_max", flags.Lookup("zoom-max"))
	viper.BindPFlag("out_scheme", flags.Lookup("out-scheme"))
	viper.BindPFlag("download_agent", flags.Lookup("agent"))

}

func main() {
	mainCommand.Execute()
}
