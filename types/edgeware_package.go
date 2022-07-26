package types

import (
	"os"
	"io"
	"path"
	"strings"
	"archive/zip"
	"encoding/json"
)

// Just.. see the Edgeware readme to explain this.
type Site struct {
	Url   string
	Query string
}

// As with this
type PromptSet struct {
	Mood      string
	Frequency int32
	Prompts   []string
}

type CaptionSet struct {
	Prefix    string
	Sentences []string
}

type EdgewarePackage struct {
	Sites []Site
	Prompts []PromptSet
	MinSentences int
	MaxSentences int
	Captions []CaptionSet
	// A path to an image
	Icon string
	// A path to an image
	Wallpaper string
	// A list of paths to audio files
	AudioFiles []string
	// A list of paths to images
	ImageFiles []string
	// A list of paths to gifs
	SubliminalFiles []string
	// A list of paths to videos
	VideoFiles []string
}

func LoadEdgewarePackage(packagePath, packageExtractDirectory string) EdgewarePackage {
	var pkg EdgewarePackage

	_, err := os.Stat(packageExtractDirectory)
	if err == nil {
		err := os.RemoveAll(packageExtractDirectory)
		if err != nil { 
			// Todo: This better
			panic(err)
		}
	} else if !os.IsNotExist(err) {
		// Todo: This better
		panic(err)
	}

	_ = os.Mkdir(packageExtractDirectory, os.ModePerm)

	zipReader, err := zip.OpenReader(packagePath)
	if err != nil {
		// Todo: This better
		panic(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := packageExtractDirectory + "/" + f.Name

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			// Todo: this better
			panic(err)
		}

		defer dstFile.Close()

		fileInArchive, err := f.Open()
		if err != nil {
			// Todo: this better
			panic(err)
		}

		defer fileInArchive.Close()
		
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			// Todo: this better
			panic(err)
		}

		if strings.HasSuffix(filePath, "web.json") {
			pkg.Sites = ParseWebJson(filePath)
		} else if strings.HasSuffix(filePath, "prompt.json") {
			pkg.Prompts, pkg.MinSentences, pkg.MaxSentences = ParsePromptJson(filePath)
		} else if strings.HasSuffix(filePath, "captions.json") {
			pkg.Captions = ParseCaptionsJson(filePath)
		} else if strings.HasSuffix(filePath, "wallpaper.png") {
			pkg.Wallpaper = filePath;
		} else {
			dir := path.Dir(filePath)

			if strings.HasSuffix(dir, "img") && (strings.HasSuffix(filePath, "png") || strings.HasSuffix(filePath, "jpg")) {
				pkg.ImageFiles = append(pkg.ImageFiles, filePath)
			} else if strings.HasSuffix(dir, "vid") && (strings.HasSuffix(filePath, "mp4") || strings.HasSuffix(filePath, "webm")) {
				pkg.VideoFiles = append(pkg.VideoFiles, filePath)
			} else if strings.HasSuffix(dir, "aud") && (strings.HasSuffix(filePath, "wav") || strings.HasSuffix(filePath, "mp3")) {
				pkg.AudioFiles = append(pkg.AudioFiles, filePath)
			} else if strings.HasSuffix(dir, "subliminals") && strings.HasSuffix(filePath, "gif") {
				pkg.SubliminalFiles = append(pkg.SubliminalFiles, filePath)
			}
		}
	}

	return pkg
}

func ParseWebJson(filePath string) []Site {
	var sites []Site

	webJsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	var raw map[string][]string
	err = json.Unmarshal(webJsonBytes, &raw)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	for i := 0; i < len(raw["urls"]); i++ {
		sites = append(sites, Site{
			Url: raw["urls"][i],
			Query: raw["args"][i],
		})
	}

	return sites
}

func ParsePromptJson(filePath string) ([]PromptSet, int, int) {
	var promptsets []PromptSet
	var min int
	var max int

	promptJsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	var raw map[string]json.RawMessage
	err = json.Unmarshal(promptJsonBytes, &raw)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	err = json.Unmarshal(raw["minLen"], &min)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}
	err = json.Unmarshal(raw["maxLen"], &max)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	var unmarshalledMoods []string
	err = json.Unmarshal(raw["moods"], &unmarshalledMoods)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}
	
	for i := 0; i < len(unmarshalledMoods); i++ {
		moodName := unmarshalledMoods[i]
		
		var unmarshalledPrompts []string
		err = json.Unmarshal(raw[moodName], &unmarshalledPrompts)
		if err != nil {
			// Todo: Handle better
			panic(err)
		}

		var freqList []int32
		err = json.Unmarshal(raw["freqList"], &freqList)
		if err != nil {
			// Todo: Handle better
			panic(err)
		}

		promptsets = append(promptsets, PromptSet{
			Mood: string(moodName),
			Frequency: freqList[i],
			Prompts: unmarshalledPrompts,
		})
	}

	return promptsets, min, max
}

func ParseCaptionsJson(filePath string) []CaptionSet {
	var captionsets []CaptionSet

	captionJsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	var raw map[string][]string
	err = json.Unmarshal(captionJsonBytes, &raw)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	for _, prefix := range raw["prefix"] {
		captionsets = append(captionsets, CaptionSet{
			Prefix: prefix,
			Sentences: raw[prefix],
		})
	}

	return captionsets
}

func ExtractString(obj json.RawMessage) string {
	var s string
	json.Unmarshal(obj, &s)
	return s
}