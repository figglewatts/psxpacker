# psxpacker
Pack textures into an atlas sort of like how it was done in PSX games. 

## Usage
```
./psxpacker <direction (width/height)> <dimension> <dir_or_file>...
```

## Features
- Direction of width makes columns of texture widths sorted by height
- Direction of height makes columns of texture heights sorted by width
- Tie break between width/height sort is filename
- Default direction is height
- Errors if all textures cannot be packed into dimension