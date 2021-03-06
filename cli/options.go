package cli

import (
	"flag"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	ImageIds []string
)

// Initialise the package.
func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.IMAGE_PATH) {
			id := strings.TrimSuffix(path.Base(asset), data.IMAGE_EXTENSION)
			ImageIds = append(ImageIds, id)
		}
	}

	sort.Sort(sort.StringSlice(ImageIds))
}

type Options struct {
	Bottom   string
	ClientId string
	Help     bool
	Image    string
	Name     string
	Top      string
}

// Parse the command line options.
func ParseOptions() Options {
	var opt Options
	var text string

	flag.BoolVar(&opt.Help, "h", false, "Show help.")
	flag.StringVar(&opt.ClientId, "cid", "", "The client id of an application registered with imgur.com. If specified, the new meme will be uploaded to imgur.com instead of being saved locally. (See README for full details.)")
	flag.StringVar(&opt.Image, "i", "", "One of the built-in templates, a URL or the path to a local file (gif, jpeg or png.) You can also use '-' to read an image from stdin.")
	flag.StringVar(&opt.Name, "o", "", "The optional name of the output file (png). If omitted, a temporary file will be used.")
	flag.StringVar(&text, "t", "", "The meme text. Separate the top and bottom banners using a pipe character.")
	flag.Parse()

	parsed := strings.Split(text, "|")
	if len(parsed) == 1 {
		opt.Top = parsed[0]
	} else {
		opt.Top = parsed[0]
		opt.Bottom = parsed[1]
	}

	return opt
}

// Validate the command line options.
func (this *Options) Valid() bool {

	if this.Image == "" {
		output.Error("An image is required")
	}

	if this.Name != "" {
		if !strings.HasSuffix(strings.ToLower(this.Name), ".png") {
			output.Error("The output file name must have the suffix of .png")
		}
	}

	if (this.Top + this.Bottom) == "" {
		output.Error("At least one piece of text is required")
	}

	return true
}

// Print the usage of this program.
func (this *Options) PrintUsage() {
	var banner string = ` _ __ ___   ___ _ __ ___   ___
| '_ ' _ \ / _ \ '_ ' _ \ / _ \
| | | | | |  __/ | | | | |  __/
|_| |_| |_|\___|_| |_| |_|\___|

`
	color.Green(banner)
	fmt.Println("")
	flag.Usage()
	fmt.Println("")

	fmt.Println("  Templates")
	fmt.Println("")
	for _, name := range ImageIds {
		color.Cyan("    " + name)
	}
	fmt.Println("")

	fmt.Println("  Examples")
	fmt.Println("")
	color.Cyan("    meme -i kirk-khan -t \"|khaaaan\"")
	color.Cyan("    meme -i brace-yourselves -t \"Brace yourselves|The memes are coming!\"")
	color.Cyan("    meme -i http://i.imgur.com/FsWetC0.jpg -t \"|China\"")
	color.Cyan("    meme -i ~/Pictures/face.png -t \"Hello\"")
	fmt.Println("")
}
