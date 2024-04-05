import customtkinter as ctk
from tkinter import Event
from PIL import Image, ExifTags
from exiftool import ExifToolHelper

class ImageInfo:
    def __init__(self, root : ctk.CTkFrame, images : list[tuple[Image.Image, str]], curr_index_var : ctk.IntVar) -> None:
        self.root = root
        self.images = images
        self.curr_index_var = curr_index_var
        
        self.root.grid_configure(ipadx=30, ipady=10)
        # self.root.grid_configure(ipadx=10, ipady=0)
        
        self.root.grid_rowconfigure([i for i in range(0, 9+2)], weight=1)
        # self.root.grid_rowconfigure((0, 1), weight=1, pad=0)
        self.root.grid_columnconfigure(0, weight=1)
        self.root.grid_columnconfigure(1, weight=2)
        self.root.grid_columnconfigure(2, weight=2)
        self.root.grid_columnconfigure(3, weight=1)
        
        self.filename_var = ctk.StringVar()
        self.cam_type_var = ctk.StringVar()
        self.lens_type_var = ctk.StringVar()
        self.date_time_var = ctk.StringVar()
        self.focal_length_var = ctk.StringVar()
        self.exposure_time_var = ctk.StringVar()
        self.aperture_var = ctk.StringVar()
        self.iso_var = ctk.StringVar()
        self.rating_var = ctk.StringVar()

        self.update_exif_data()
        self.create_widgets()
    
    def create_widgets(self):
        line_height = 1
        font = ("Arial", 18)
        # filename
        self.filename_label = ctk.CTkLabel(self.root, textvariable=self.filename_var, height=line_height, pady=0, font=font)
        self.filename_label.grid(row=1, column=1, columnspan=2, sticky="w")
        
        # camera type
        self.cam_type_label = ctk.CTkLabel(self.root, textvariable=self.cam_type_var, height=line_height, pady=0, font=font)
        self.cam_type_label.grid(row=2, column=1, columnspan=2, sticky="w")
        
        # lens type
        self.lens_type_label = ctk.CTkLabel(self.root, textvariable=self.lens_type_var, height=line_height, pady=0, font=font)
        self.lens_type_label.grid(row=3, column=1, columnspan=2, sticky="w")
        
        # date and time
        self.date_time_label = ctk.CTkLabel(self.root, textvariable=self.date_time_var, height=line_height, pady=0, font=font)
        self.date_time_label.grid(row=4, column=1, columnspan=2, sticky="w")

        # focal length
        self.focal_length_label = ctk.CTkLabel(self.root, text="Focal length:", height=line_height, pady=0, font=font)
        self.focal_length_label.grid(row=5, column=1, sticky="w")
        self.focal_length_label_val = ctk.CTkLabel(self.root, textvariable=self.focal_length_var, height=line_height, pady=0, font=font)
        self.focal_length_label_val.grid(row=5, column=2, sticky="e")
        
        # exposure time
        self.exposure_time_label = ctk.CTkLabel(self.root, text="Exposure time:", height=line_height, pady=0, font=font)
        self.exposure_time_label.grid(row=6, column=1, sticky="w")
        self.exposure_time_label_val = ctk.CTkLabel(self.root, textvariable=self.exposure_time_var, height=line_height, pady=0, font=font)
        self.exposure_time_label_val.grid(row=6, column=2, sticky="e")
        
        # aperture value
        self.aperture_label = ctk.CTkLabel(self.root, text="Aperture:", height=line_height, pady=0, font=font)
        self.aperture_label.grid(row=7, column=1, sticky="w")
        self.aperture_label_val = ctk.CTkLabel(self.root, textvariable=self.aperture_var, height=line_height, pady=0, font=font)
        self.aperture_label_val.grid(row=7, column=2, sticky="e")
        
        # ISO
        self.iso_label = ctk.CTkLabel(self.root, text="ISO:", height=line_height, font=font)
        self.iso_label.grid(row=8, column=1, sticky="w")
        self.iso_label_val = ctk.CTkLabel(self.root, textvariable=self.iso_var, height=line_height, font=font)
        self.iso_label_val.grid(row=8, column=2, sticky="e")
        
        font = (font[0], font[1] + 6)
        # rating
        self.rating_label = ctk.CTkLabel(self.root, textvariable=self.rating_var, height=line_height, font=font)
        self.rating_label.grid(row=9, column=1, columnspan=2, sticky="nsew")
    
    def reset_vars(self):
        curr_image_path = self.images[self.curr_index_var.get()][1]
        if curr_image_path != None:
            self.filename_var.set(curr_image_path.split("/")[-1])
        else:
            self.filename_var.set("")
            
        self.cam_type_var.set("")
        self.lens_type_var.set("")
        self.date_time_var.set("")
        self.focal_length_var.set("")
        self.exposure_time_var.set("")
        self.aperture_var.set("")
        self.iso_var.set("")
        self.rating_var.set("")
    
    def update_exif_data(self):
        curr_image_path = self.images[self.curr_index_var.get()][1]
        
        with ExifToolHelper() as et:
            exif_info = et.get_tags(
                [curr_image_path],
                [
                    "Model",
                    "LensModel",
                    "DateTimeOriginal",
                    "FocalLength",
                    "ExposureTime",
                    "FNumber",
                    "ISO",
                    "Rating"
                ])[0]
        
        if exif_info is None:
            self.reset_vars()
            return
        
        if curr_image_path != None:
            self.filename_var.set(curr_image_path.split("/")[-1])
        else:
            self.filename_var.set("")
        
        if "EXIF:Model" in exif_info.keys():
            self.cam_type_var.set(exif_info["EXIF:Model"])
        else:
            self.cam_type_var.set("Unknown")
        
        if "EXIF:LensModel" in exif_info.keys():
            self.lens_type_var.set(exif_info["EXIF:LensModel"])
        else:
            self.lens_type_var.set("Unknown")

        if "EXIF:DateTimeOriginal" in exif_info.keys():
            self.date_time_var.set(exif_info["EXIF:DateTimeOriginal"])
        else:
            self.date_time_var.set("Unknown")
        
        if "EXIF:FocalLength" in exif_info.keys():
            self.focal_length_var.set(f"{exif_info["EXIF:FocalLength"]}mm")
        else:
            self.focal_length_var.set("Unknown")
        
        if "EXIF:ExposureTime" in exif_info.keys():
            self.exposure_time_var.set(f"1/{1/exif_info["EXIF:ExposureTime"]}s")
        else:
            self.exposure_time_var.set("Unknown")
        
        if "EXIF:FNumber" in exif_info.keys():
            self.aperture_var.set(f"F{exif_info["EXIF:FNumber"]}")
        else:
            self.aperture_var.set("Unknown")
        
        if "EXIF:ISO" in exif_info.keys():
            self.iso_var.set(exif_info["EXIF:ISO"])
        else:
            self.iso_var.set("Unknown")
        
        if "XMP:Rating" in exif_info.keys():
            rating = exif_info["XMP:Rating"]
            rating_str = ""
            for i in range(1, 6):
                if i <= rating:
                    rating_str += "★"
                else:
                    rating_str += "☆"
            self.rating_var.set(rating_str)
        else:
            self.rating_var.set("Rating Unknown")
        

if __name__ == "__main__":
    from imageApp import ImageApp
    
    root = ctk.CTk()
    root.title("Image Viewer")
    app = ImageApp(root, "~/Pictures")
    app.create_widgets()
    root.mainloop()
