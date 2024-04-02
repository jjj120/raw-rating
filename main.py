# /bin/python3
# -*- coding: utf-8 -*-

import os
from PIL import Image, ImageTk, ExifTags
import tkinter as tk
from itertools import cycle
from imageDisplay import ImageDisplay


def load_images_from_folder(folder: str):
    images = []
    for filename in os.listdir(folder):
        if filename.endswith('.jpg') or filename.endswith('.jpeg') or filename.endswith('.JPG') or filename.endswith('.JPEG'):
            img_path = os.path.join(folder, filename)
            img = Image.open(img_path)
            images.append((img, img_path))
    return images

def main():
    folder = '/mnt/daten/1 Fotos/2024/2024-04-01 Ostermontag'
    images = load_images_from_folder(folder)

    if images:
        root = tk.Tk()
        app = ImageDisplay(root, images)
        root.mainloop()
    else:
        print("No CR3 images found in the folder.")

if __name__ == "__main__":
    main()
