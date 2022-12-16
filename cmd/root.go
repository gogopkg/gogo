package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gogopkg/gogo/pkg/gogo"
)

var rootCmd = &cobra.Command{
	Use: "Rune",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := serve(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("gogofile", "f", "Gogofile.yaml", "Gogofile name")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func serve(ctx context.Context) error {
	gogofile := viper.GetString("gogofile")

	b, err := ioutil.ReadFile(gogofile)
	if err != nil {
		return errors.Wrapf(err, "read %v", gogofile)
	}

	gg := gogo.NewGogo()
	if err := gg.LoadGlobal(ctx, b); err != nil {
		return errors.Wrapf(err, "load")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(gg.GetData()); err != nil {
		return errors.Wrapf(err, "json encode")
	}

	return nil
}
