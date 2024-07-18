# RAW Rating tool

## Description

This is a tool to rate RAW image files and write these ratings directly to the exif data of the RAW and JPG file. It is useful for photographers who want to rate their images before importing them into a photo management software (Lightroom, Darktable, etc.).

The tool displays the image and allows the user to rate it with a number from 0 to 5. The rating is written to the exif data of the RAW and JPG file. The tool is written in Go and uses the gotk3 library for the GTK3 based GUI. The tool is still in development and has some bugs and performance issues. (See TODO section for more information)

It needs the exiftool to be installed on the system. The exiftool is used to write the rating to the exif data of the RAW and JPG file.

## Installation

1. Install the exiftool from [exiftool.org](https://exiftool.org/)
2. Install Go from [go.dev](https://go.dev/)
3. Clone the repository
4. Run `go build` in the repository folder
5. Run the resulting binary

## Usage


1. Run the tool
1. Select the folder containing the RAW files via the top menu
1. Cycle through the images
1. Rate the image with the number keys (0-5)
1. To see details, click on the image
1. Press `Esc` or `q` to quit the tool

## Development

This tool was developed with `Go 1.22.4` and tested on Fedora 40. It should work on other GNOME based Linux distros as well, but I haven't tested it.

If you want to contribute, feel free to fork the repository and create a pull request. The Tool should also work for other RAW formats as long as it is supported by exiftool, but I haven't tested it. If you encounter any bugs or have feature requests, feel free to open an issue.

If you want to get more debugging information, you can set the `DEBUG` or `INFO` environment variable to `"true"` before running the tool.

```bash
export DEBUG="true"
./raw-rating <path-to-image-folder>
```

```bash
export INFO="true"
./raw-rating <path-to-image-folder>
```


## TODO

- [ ] Dynamically loaded image preview in the list
- [ ] Add rating to the list items
- [ ] Buttons for rating
- [ ] Settings for the tool
- [ ] Customizable shortcuts
- [ ] **Performance improvements**
- [ ] Bugfixes

## Dependencies

- [gotk3](https://github.com/gotk3/gotk3)
- [exiftool](https://exiftool.org/)
- [go-exiftool](https://github.com/barasher/go-exiftool)
- [logrus](https://github.com/sirupsen/logrus)
- [Go](https://go.dev/)

## Supported RAW-Formats

The tool should work with all RAW formats supported by exiftool. They are listed in the [exiftool documentation](https://exiftool.org/#supported).
