# psxpacker
Pack textures into an atlas sort of like how it was done in PSX games. 

## Prerequisites
- Go 1.17+

## Installation
```
$ go get github.com/Figglewatts/psxpacker/cmd/psxpacker
```

## Usage
```
./psxpacker direction width height output_path dir_or_file...

Pack textures into an atlas sort of like how it was done in PSX games.

Arguments:
    direction    The direction to pack images along. One of ['width', 'height'].
    width        The width of the packed image to produce.
    height       The height of the packed image to produce.
    output_path  The output path of the packed image.
    dir_or_file  A PNG image or directory containing PNG images to pack.
```

## Features
- Direction of width makes columns of texture widths sorted by height
- Direction of height makes columns of texture heights sorted by width
- Tie break between width/height sort is filename
- Default direction is height
- Errors if all textures cannot be packed into dimension