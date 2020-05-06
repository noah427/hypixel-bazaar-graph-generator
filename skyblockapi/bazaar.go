package skyblockapi

import (
	"fmt"

	

	"github.com/imroc/req"
)

// GetItemStats : get's the stats of a hypixel bazaar item when given it's ID
func GetItemStats(id string, key string) string {
	response, err := req.Get(fmt.Sprintf("https://api.hypixel.net/skyblock/bazaar/product?key=%s&productId=%s", key, id))

	if err != nil {
		fmt.Println(err)
	}

	return response.String()
}
