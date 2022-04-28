package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var Verbose bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cobra-viper-sample",
	Short: "Sample application using cobra",
	Long: `This application is simple program to learn spf13/cobra
spf13/cobra looks really nice CLI framework.
I want to create lovely CLI program with this framework :)
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra-viper-sample.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	// Cobra also supports local flags, which will only run
	RootCmd.PersistentFlags().Bool("logging-prod", true, "")
	RootCmd.PersistentFlags().StringP("logging-level", "l", "info", "")
	RootCmd.PersistentFlags().StringP("beat-saber-path", "p", "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/", "")
	RootCmd.PersistentFlags().Float32P("min-rating-percent", "r", .5, "")
	RootCmd.PersistentFlags().IntP("workers", "w", 4, "")
	RootCmd.PersistentFlags().IntP("download-delay", "d", 5, "")
	RootCmd.PersistentFlags().IntP("song-download-amount", "a", 20, "")
	RootCmd.PersistentFlags().StringP("user_name", "u", "", "")
	// when this action is called directly.
	songsInit()
	songsSearchInit()
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	//log.Fatal(err)
	//}
	//
	//viper.SetConfigName(".cobra-viper-sample") // name of config file (without extension)
	//if cfgFile != "" {                         // enable ability to specify config file via flag
	//	fmt.Println(">>> cfgFile: ", cfgFile)
	//	viper.SetConfigFile(cfgFile)
	//	configDir := path.Dir(cfgFile)
	//	if configDir != "." && configDir != dir {
	//		viper.AddConfigPath(configDir)
	//	}
	//}
	//
	//viper.AddConfigPath(dir)
	//viper.AddConfigPath(".")
	//viper.AddConfigPath("$HOME")
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//} else {
	//	fmt.Println(err)
	//}
	//viper.WatchConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})
}
