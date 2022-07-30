package Goonware

import (
	"fmt"
	"strings"

	types "goonware/types"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

var selectedTag int32
var newTag string

func DriveFillterTab(c *types.Config) []g.Widget {
	largerFont := g.GetDefaultFonts()[0].SetSize(20)

	return []g.Widget{
		g.PopupModal("Add New Tag").Layout(g.Layout{
			g.InputText(&newTag).Size(300),
			g.Button("Ok").OnClick(func() {
				if !(strings.TrimSpace(newTag) == "") {
					c.DriveFillerTags = append(c.DriveFillerTags, newTag)
					newTag = ""
				}
				g.CloseCurrentPopup()
			}),
		}),

		g.Checkbox("On", &c.DriveFiller),

		ConditionOrNothing(c.DriveFiller, g.Layout{
			g.Row(
				LabelSliderTooltip("Fill delay", &c.DriveFillerDelay, 10, 3000, 200,
					"How many milliseconds to wait before writing another image",
					FormatMillisecondSlider),
			),
			g.Row(
				g.Button("Select base").OnClick(func() { SelectBase(c) }),
				g.Label("("+c.DriveFillerBase+")"),
			),
			StandardSeparation(),

			g.Row(
				g.Label("Image source"),
				g.RadioButton("Use Package", c.DriveFillerImageSource == 0).
					OnChange(func() { c.DriveFillerImageSource = 0 }),
				g.Tooltip("Fill the drive with files in the currently loaded package's img directory"),

				g.RadioButton("Download", c.DriveFillerImageSource == 1).
					OnChange(func() { c.DriveFillerImageSource = 1 }),
				g.Tooltip("Fill the drive with images downloaded from a booru of your choosing"),
			),
			StandardSeparation(),

			ConditionOrNothing(c.DriveFillerImageSource == 1, g.Layout{
				g.Row(g.Label("Image Downloader").Font(largerFont)),
				StandardSeparation(),

				g.Row(
					g.Label("Booru"),
					g.Label("https://e621.net/"),
				),

				g.Row(
					g.RadioButton("Anything", !c.DriveFillerImageUseTags).
						OnChange(func() { c.DriveFillerImageUseTags = false }),
					g.RadioButton("Specific Tags", c.DriveFillerImageUseTags).
						OnChange(func() { c.DriveFillerImageUseTags = true }),
				),
				g.Child().Layout(TagsLayout(c)).Size(300, 300),
				g.Row(
					g.Button("+").OnClick(func() { g.OpenPopup("Add New Tag") }),
					g.Button("-").OnClick(func() {
						c.DriveFillerTags = RemoveElement(c.DriveFillerTags, selectedTag)
					}),
				),

				g.Row(
					g.Checkbox("Minimum score", &c.DriveFillerDownloadMinimumScoreToggle),
					ConditionOrNothing(c.DriveFillerDownloadMinimumScoreToggle,
						g.Layout{g.SliderInt(&c.DriveFillerDownloadMinimumScoreThreshold, -50, 100).
							Size(150)},
					),
				),
				StandardSeparation(),
			}),
		}),
	}
}

func SelectBase(c *types.Config) {
	folder, err := dialog.Directory().Browse()

	if err != nil && err != dialog.Cancelled {
		dialog.Message("%s", fmt.Sprintf("Couldn't select directory: %s", err.Error())).Error()
	} else if err == nil {
		c.DriveFillerBase = folder
	}
}

func TagsLayout(c *types.Config) g.Layout {
	var r g.Layout

	for idx, tag := range c.DriveFillerTags {
		r = append(r, g.Layout{
			g.Selectable(tag).
				OnClick(func() { selectedTag = int32(idx) }).
				Selected(selectedTag == int32(idx)),
		})
	}

	return r
}

func RemoveElement(slice []string, idx int32) []string {
	return append(slice[:idx], slice[idx+1:]...)
}
