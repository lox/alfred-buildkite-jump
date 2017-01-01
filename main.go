package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/pascalw/go-alfred"
)

func main() {
	app := kingpin.New("alfred-buildkite-jump", "Jump to buildkite pipelines")
	configurePipelinesCommand(app)
	configureUpdateCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func alfredError(err error) *alfred.AlfredResponseItem {
	return &alfred.AlfredResponseItem{
		Valid:    false,
		Uid:      "error",
		Title:    "Error Occurred",
		Subtitle: err.Error(),
		Icon:     "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/AlertStopIcon.icns",
	}
}
