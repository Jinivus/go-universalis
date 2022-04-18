package universalis

import (
	"context"
	"fmt"
	"github.com/jinivus/go-xivapi"
	"os"
)

func FindXIVApiItem(query string) *xivapi.Item {
	fmt.Printf("Searching xivapi for %s...\n", query)
	xivapiClient := xivapi.NewClient(nil)
	matches, _, err := xivapiClient.Search.Items(context.Background(), query)
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

	return matches.Items[0]
}
