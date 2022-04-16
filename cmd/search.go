/*
Copyright © 2022 Matt Cooper <jinivus@jinivus.com>

*/
package cmd

import (
	"context"
	"fmt"
	universalis "github.com/jinivus/go-universalis/pkg"
	"github.com/jinivus/go-xivapi"
	"github.com/spf13/cobra"
	"os"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Searching xivapi for %s...\n", args[0])
		xivapiClient := xivapi.NewClient(nil)
		matches, _, err := xivapiClient.Search.Items(context.Background(), args[0])
		if err != nil {
			fmt.Printf("error searching xivapi: %s\n", err)
			os.Exit(1)
		}
		if len(matches.Items) > 1 {
			fmt.Printf("more than one match found\n")
			os.Exit(1)
		}
		if len(matches.Items) == 0 {
			fmt.Printf("no match found\n")
			os.Exit(1)
		}
		fmt.Printf("Item found on xivapi, item id is %d...\n", matches.Items[0].ID)

		server, _ := cmd.Flags().GetString("server")
		hq, _ := cmd.Flags().GetBool("hq")

		universalisClient := universalis.NewClient(nil)
		listingResult, _, err := universalisClient.Listings.ListingsWithOptions(context.Background(), server, fmt.Sprintf("%d", matches.Items[0].ID), universalis.ListingOptions{HQOnly: hq})
		if err != nil {
			fmt.Printf("error gettings universalis listings: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Lowest current price on %s: %d gil", server, listingResult.MinPriceHQ)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	searchCmd.PersistentFlags().String("server", "ravana", "Server to search for items on, by name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	searchCmd.Flags().BoolP("hq", "q", false, "Search hq only")
}
