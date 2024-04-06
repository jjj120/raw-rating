# CR3 Rating tool

## Description

This is a tool to rate CR3 files (Canon RAW 3) and write these ratings directly to the exif data of the CR3 and JPG file. It is useful for photographers who want to rate their images before importing them into a photo management software (Lightroom, Darktable, etc.). 

The tool displays the image and allows the user to rate it with a number from 0 to 5. The rating is written to the exif data of the CR3 and JPG file. The tool is written in Python and uses the [Customtkinter](https://customtkinter.tomschimansky.com/) library for the GUI. The tool is still in development and has some bugs and performance issues. (See TODO section for more information)

It needs [exiftool](https://exiftool.org/) to be installed on the system. You also have to shoot in RAW+JPG mode, because the JPG file is used to display the image. The rating is written to the CR3 and JPG file.

## Installation

1. Install exiftool from [https://exiftool.org/](https://exiftool.org/)
2. Install the required packages with `pip install -r requirements.txt`
3. Run the tool with `python main.py`

## Usage

1. Run the tool
2. Select the folder containing the CR3 files
3. Cycle through the images with the arrow keys
4. Rate the image with the number keys (`0-5`)
5. Zoom in/out with the `mouse wheel`
6. Reset zoom with `double click`
7. Press `Esc` to quit the tool

## Images

Path Selector:

<img src=".images/path_selector.png" alt="Select folder" width="30%"/>

Image Display:

<img src=".images/image_display.png" alt="Image Display" width="100%"/>

## Development

This tool was developed with Python 3.12.2 and tested on Fedora 39. It should work on other Linux distros and Windows as well, but I haven't tested it.

**If you want to contribute, feel free to fork the repository and create a pull request.**
The Tool should also work for other RAW formats as long as it is supported by exiftool, but I haven't tested it. If you want to use other RAW formats, you have to change the file extension in the code, because it is currently hardcoded to CR3.

### TODO
- [ ] Support for other RAW formats
- [ ] Dynamically loaded image preview in the list
- [ ] Add rating to the list
- [ ] Full screen Gallery mode
- [ ] Settings menu
- [ ] Customizable key bindings
- [ ] Deleting images key binding
- [ ] Buttons for all functions
- [ ] Bug fixes
- [ ] **Performance improvements**

## Dependencies
- [exiftool](https://exiftool.org/)
- [Pillow](https://pillow.readthedocs.io/en/stable/)
- [Customtkinter](https://customtkinter.tomschimansky.com/)
- [CTkListbox](https://github.com/Akascape/CTkListbox)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.