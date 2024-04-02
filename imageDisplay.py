import os
from PIL import Image, ImageTk
import tkinter as tk
from exiftool import ExifToolHelper


class ImageDisplay:
    def __init__(self, master: tk.Tk, images: list):
        self.master : tk.Tk = master
        self.images = images
        self.current_index = 0
        self.photo_image = None

        self.canvas = tk.Canvas(
            master,
            background='#3B3B3B',
            borderwidth=0,
            highlightthickness=0,
            width=1920,
            height=1080
        )
        self.canvas.pack(fill=tk.BOTH, expand=tk.YES)

        # Bind keyboard events
        master.bind('<Right>', lambda event: self.next_image())
        master.bind('<Left>', lambda event: self.prev_image())
        master.bind('<Escape>', lambda event: master.quit())
        master.bind('<Configure>', lambda event: self.resizing(event))

        master.bind('0', lambda event: self.change_rating(0))
        master.bind('1', lambda event: self.change_rating(1))
        master.bind('2', lambda event: self.change_rating(2))
        master.bind('3', lambda event: self.change_rating(3))
        master.bind('4', lambda event: self.change_rating(4))
        master.bind('5', lambda event: self.change_rating(5))

        self.display_image()

    def change_rating(self, rating: int):
        with ExifToolHelper() as et:
            et.set_tags(
                [self.images[self.current_index][1], self.images[self.current_index][1].replace(".JPG", ".CR3")],
                tags={"Rating": rating},
                params=["-P", "-overwrite_original"]
            )
        print(f"Rating of {self.images[self.current_index][1]} changed to {rating}")
        self.images[self.current_index][0].close()
        self.images[self.current_index] = (Image.open(self.images[self.current_index][1]), self.images[self.current_index][1])
        self.display_image()

    def resizing(self, event):
        if event.width > 1 and event.height > 1:
            self.display_image()

    def display_image(self):
        if self.photo_image:
            self.canvas.delete(tk.ALL)

        self.master.update()
        height = self.master.winfo_height()
        width = self.master.winfo_width()

        img, img_path = self.images[self.current_index]
        img = img.copy()
        img.thumbnail((width, height), Image.LANCZOS)

        self.photo_image = ImageTk.PhotoImage(img)
        self.canvas.create_image(width / 2, height / 2, anchor=tk.CENTER, image=self.photo_image)

        # Display EXIF information
        self.update_exif_info(img_path)

    def next_image(self):
        self.current_index = (self.current_index + 1)
        self.check_index()
        self.display_image()

    def prev_image(self):
        self.current_index = (self.current_index - 1)
        self.check_index()
        self.display_image()

    def check_index(self):
        if self.current_index >= len(self.images):
            self.current_index = self.current_index - 1
        elif self.current_index < 0:
            self.current_index = 0

    def get_exif_info(self, img_path):
        img = Image.open(img_path)
        exif_data = img._getexif()
        if exif_data is None:
            return "No EXIF data found."

        filename : str = os.path.basename(img_path)
        exif_info = f"{filename.replace(".JPG", "")}\n"
        exif_info += f"{exif_data.get(272)}\n"
        exif_info += f"{exif_data.get(306)}\n"
        exif_info += f"Focal length:  {exif_data.get(37386)}mm\n"
        exif_info += f"Exposure time: 1/{1/exif_data.get(33434)}s\n"
        exif_info += f"Aperture:      F{exif_data.get(33437)}\n"
        exif_info += f"ISO:           {exif_data.get(34855)}\n"

        with ExifToolHelper() as et:
            rating = et.get_tags([img_path, img_path.replace(".JPG", ".CR3")], ["Rating"])[0]["XMP:Rating"]

        exif_info += f"Rating:        {rating}\n"

        return exif_info

    def update_exif_info(self, img_path):
        font_height = 12
        exif_info = self.get_exif_info(img_path)
        x1, y1 = 10, 10
        x2, y2 = 270, font_height*2*exif_info.count("\n")+x1+10

        self.create_rectangle(x1, y1, x2, y2, fill='black', alpha=0.3)
        self.canvas.create_text(10+x1, 10+y1, anchor='nw', text=exif_info, fill='white', font=('DejaVu Sans Mono', font_height))

    def create_rectangle(self, x1, y1, x2, y2, **kwargs):
        if 'alpha' in kwargs:
            alpha = int(kwargs.pop('alpha') * 255)
            fill = kwargs.pop('fill')
            fill = self.master.winfo_rgb(fill) + (alpha,)
            image = Image.new('RGBA', (x2-x1, y2-y1), fill)
            self.exif_info_bg = ImageTk.PhotoImage(image)
            self.canvas.create_image(x1, y1, image=self.exif_info_bg, anchor='nw')
        self.canvas.create_rectangle(x1, y1, x2, y2, **kwargs)


if __name__ == "__main__":
    import main
    main.main()
