package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Figglewatts/psxpacker/pkg/psxpacker"
)

type Args struct {
	Direction          psxpacker.Direction
	Width              int
	Height             int
	Output             string
	DirectoriesOrFiles []string
}

var (
	//lint:ignore ST1012 This behaviour is in the go std library
	ExitProgram = errors.New("exit program")
)

const (
	UsageText       = "./psxpacker direction width height output_path dir_or_file..."
	DescriptionText = "Pack textures into an atlas sort of like how it was done in PSX games. "
	MinArgs         = 5
	HelpText        = `Arguments:
    direction    The direction to pack images along. One of ['width', 'height'].
    width        The width of the packed image to produce.
	height       The height of the packed image to produce.
    output_path  The output path of the packed image.
    dir_or_file  A PNG image or directory containing PNG images to pack.
`
)

func parseArgs(args []string) (Args, error) {
	if len(args) >= 1 && (args[0] == "-h" || args[0] == "--help") {
		// print help text and exit
		fmt.Println(UsageText)
		fmt.Println()
		fmt.Println(DescriptionText)
		fmt.Println()
		fmt.Print(HelpText)
		fmt.Println()
		return Args{}, ExitProgram
	}

	// ensure correct number of args parsed
	if len(args) < MinArgs {
		return Args{}, fmt.Errorf("at least %v arguments required\nUsage: %s", MinArgs, UsageText)
	}

	// attempt to parse direction
	direction, err := psxpacker.DirectionFromString(args[0])
	if err != nil {
		return Args{}, fmt.Errorf("direction: %v\nUsage: %s", err, UsageText)
	}

	// attempt to parse width
	width, err := strconv.Atoi(args[1])
	if err != nil {
		return Args{}, fmt.Errorf("width: %v\nUsage: %s", err, UsageText)
	}

	// attempt to parse height
	height, err := strconv.Atoi(args[2])
	if err != nil {
		return Args{}, fmt.Errorf("height: %v\nUsage: %s", err, UsageText)
	}

	return Args{
		direction,
		width,
		height,
		args[3],
		args[4:],
	}, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func loadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening image: %v", err)
	}
	defer reader.Close()

	loadedImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("reading image: %v", err)
	}

	return loadedImage, nil
}

func loadInputImages(directoriesOrFiles []string) ([]psxpacker.InputImage, error) {
	// first gather all images we need to load
	imagesToLoad := []string{}
	for _, directoryOrFile := range directoriesOrFiles {
		isDir, err := isDirectory(directoryOrFile)
		if err != nil {
			return nil, err
		}

		if isDir {
			// walk the directory tree for png files to load
			filepath.WalkDir(directoryOrFile, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !d.IsDir() && strings.HasSuffix(path, ".png") {
					imagesToLoad = append(imagesToLoad, path)
				}
				return nil
			})
		} else if !strings.HasSuffix(directoryOrFile, ".png") {
			// if the file isn't a PNG then skip it
			fmt.Fprintf(os.Stderr, "skipping file '%v' - not a png image", directoryOrFile)
		} else {
			// the file is a PNG, we need to load it
			imagesToLoad = append(imagesToLoad, directoryOrFile)
		}
	}

	// now load the images
	inputImages := []psxpacker.InputImage{}
	for _, imageToLoad := range imagesToLoad {
		loadedImage, err := loadImage(imageToLoad)
		if err != nil {
			return nil, fmt.Errorf("image at '%v': %v", imageToLoad, err)
		}

		fileName := filepath.Base(imageToLoad)
		imageName := fileName[:len(fileName)-len(filepath.Ext(fileName))]

		inputImage := psxpacker.InputImage{
			Name:  imageName,
			Image: loadedImage,
		}
		inputImages = append(inputImages, inputImage)
	}

	return inputImages, nil
}

func writePackResult(packResult psxpacker.PackResult, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("opening output file: %v", err)
	}
	defer f.Close()

	if err = png.Encode(f, packResult.Image); err != nil {
		return fmt.Errorf("writing packed image: %v", err)
	}
	return nil
}

func run() error {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		if err == ExitProgram {
			return err
		}

		return fmt.Errorf("parsing args: %v", err)
	}

	inputImages, err := loadInputImages(args.DirectoriesOrFiles)
	if err != nil {
		return fmt.Errorf("loading images: %v", err)
	}

	packResult, err := psxpacker.Pack(args.Direction, args.Width, args.Height, inputImages)
	if err != nil {
		return fmt.Errorf("packing images: %v", err)
	}

	err = writePackResult(packResult, args.Output)
	if err != nil {
		return fmt.Errorf("writing output: %v", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		if err == ExitProgram {
			os.Exit(0)
			return
		}

		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
