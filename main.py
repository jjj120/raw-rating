# /bin/python3
# -*- coding: utf-8 -*-

import os
from PIL import Image
import customtkinter as ctk
from imageApp import ImageApp
from util import check_image_folder

def main():
    folder = '/mnt/daten/1 Fotos/2024/2024-04-01 Ostermontag'
    

    if check_image_folder(folder):
        root = ctk.CTk()
        root.title("CR3 Image Viewer")
        # app = ImageDisplay(root, images)
        app = ImageApp(root, folder)
        root.mainloop()
    else:
        print("No CR3 images found in the folder.")

if __name__ == "__main__":
    main()
