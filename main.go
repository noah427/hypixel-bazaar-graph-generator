package main

import (
	"bytes"

	"github.com/wcharczuk/go-chart"

	"github.com/tidwall/gjson"

	"image"

	"io/ioutil"

	"os"

	g "github.com/AllenDang/giu"

	"fmt"

	SKapi "./skyblockapi"
)

var (
	DisplayImage     bool
	ItemID           string
	StatMeasured     string
	ApiKey           string
	CurrentImageName string
	LoadedTexture    *g.Texture
)

func main() {
	CurrentImageName = ""

	DisplayImage = false

	window := g.NewMasterWindow("Graph Display", 700, 200, g.MasterWindowFlagsNotResizable, nil)

	window.Main(displayLoop)

}

func displayLoop() {
	g.SingleWindow("Main", g.Layout{
		g.InputText("Hypixel Api Key", float32(150), &ApiKey),
		g.InputText("Hypixel Bazaar Item ID (ALL UPPERCASE)", float32(100), &ItemID),
		g.InputText("What Statistic of the item do you want measured (options below) ", float32(150), &StatMeasured),

		g.Line(
			g.Button("Generate-Graph", generateBuyGraph),
		),

		// g.Custom(func() {
		// 	if !DisplayImage {
		// 		g.Label("No Image Generated during this session").Build()
		// 	} else {
		// 		fmt.Println("Image Load Attempt")
		// 		g.Image(LoadedTexture, 600, 400).Build()
		// 	}
		// }),
	})
}

func generateBuyGraph() {

	json := SKapi.GetItemStats(ItemID, ApiKey)

	xValues, yValues := parseJSON(json, StatMeasured)

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    fmt.Sprintf("%s | %s", ItemID, StatMeasured),
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	file, err := os.Create(fmt.Sprintf("%s | %s.png", ItemID, StatMeasured))

	if err != nil {
		fmt.Println(err)
	}

	n, err := file.Write(buffer.Bytes())

	fmt.Printf("%d bytes written\n", n)

	file.Close()

	CurrentImageName = fmt.Sprintf("%s | %s.png", ItemID, StatMeasured)

	LoadedTexture = cImageToTexture()

	DisplayImage = true

	g.Update()

}

func parseJSON(json string, stat string) ([]float64, []float64) {

	output := gjson.Get(json, "product_info.week_historic").Array()
	fmt.Println("history-len == ", len(output))

	xValues := make([]float64, len(output))

	for i := range xValues {
		xValues[i] = float64(i)
	}

	yValues := make([]float64, len(output))

	for i := range yValues {
		yValues[i] = float64(gjson.Get(json, fmt.Sprintf("product_info.week_historic.%d.%s", i, stat)).Num)
	}

	return xValues, yValues

}

