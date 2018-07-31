package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/bols-blue-org/spoon_cast_downloader/src/spoon/cast"
)

func main() {
	err := ui.Main(func() {
		input := ui.NewEntry()
		button := ui.NewButton("Download")
		greeting := ui.NewLabel("")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Enter Cast ID:"), false)
		box.Append(input, false)
		box.Append(button, false)
		box.Append(greeting, false)
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			data, err := cast.LoadMetaData(input.Text())
			if err != nil {
				greeting.SetText("Error, " + err.Error() + "!")
			}
			fmt.Printf("%v", data)
			err = data.Download()
			if err != nil {
				greeting.SetText("Error, " + err.Error() + "!")
			}

			greeting.SetText("Download, " + input.Text() + "!")
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
