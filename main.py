# /bin/python3
# -*- coding: utf-8 -*-

import os
from dotenv import dotenv_values
from PIL import Image
import customtkinter as ctk
from imageApp import ImageApp
from util import check_image_folder

def main():
    dotenv = dotenv_values(".env")
    
    folder = dotenv.get("START_FOLDER")
    
    if check_image_folder(folder):
        root = ctk.CTk()
        app = ImageApp(root, folder)
        app.create_widgets()
        root.mainloop()
    else:
        print("No CR3 images found in the folder.")

if __name__ == "__main__":
    main()
