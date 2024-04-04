import customtkinter as ctk
from tkinter import Event
from PIL import Image, ImageTk
from imageInfo import ImageInfo
from imageList import ImageList
from imageView import ImageView
from exiftool import ExifToolHelper

class ImageDisplay:
    def __init__(self, root: ctk.CTkFrame, images: list, changeToPathSelection : callable = None, changeToImageList : callable = None) -> None:
        self.root = root
        self.root.grid_rowconfigure(0, weight=2)
        self.root.grid_rowconfigure(1, weight=15)
        self.root.grid_rowconfigure(2, weight=0)
        self.root.grid_columnconfigure(0, weight=1)
        self.root.grid_columnconfigure(1, weight=5)
        
        self.images = images
        self.changeToPathSelection = changeToPathSelection if changeToPathSelection else lambda: None
        self.changeToImageList = changeToImageList if changeToImageList else lambda: None
        
        self.current_index_var = ctk.IntVar()
        
        self.create_widgets()

    def next_image(self):
        self.current_index_var.set(min(self.current_index_var.get() + 1, len(self.images) - 1))
        self.refresh_image_list()
        self.refresh_image_view()
        self.refresh_image_info()
    
    def previous_image(self):
        self.current_index_var.set(max(self.current_index_var.get() - 1, 0))
        self.refresh_image_list()
        self.refresh_image_view()
        self.refresh_image_info()

    def create_widgets(self):
        self.padding = 5
        self.image_info_frame = ctk.CTkFrame(self.root)
        self.image_info_frame.grid(row=0, column=0, sticky="nsew", padx=self.padding, pady=self.padding)
        
        self.control_frame = ctk.CTkFrame(self.root)
        self.control_frame.grid(row=2, column=0, sticky="nsew", padx=self.padding, pady=self.padding)
        
        self.image_view_frame = ctk.CTkFrame(self.root)
        self.image_view_frame.grid(row=0, column=1, rowspan=3, sticky="nsew", padx=self.padding, pady=self.padding)
        
        self.image_list_frame = ctk.CTkFrame(self.root)
        self.image_list_frame.grid(row=1, column=0, sticky="nsew", padx=self.padding, pady=self.padding)
        
        self.create_image_info()
        self.create_image_list()
        self.create_image_view()
        self.create_control()
    
    def create_image_info(self):
        self.image_info = ImageInfo(self.image_info_frame, self.images, self.current_index_var)
    
    def create_image_list(self):
        self.image_list = ImageList(self.image_list_frame, self.images, self.current_index_var, self.refresh_new_image_selection)
        
    def create_image_view(self):
        self.image_view = ImageView(self.image_view_frame, self.images, self.current_index_var)
    
    def refresh_image_info(self):
        self.image_info.update_exif_data()
    
    def refresh_image_list(self):
        self.image_list.refresh_metadata()
    
    def refresh_image_view(self):
        self.image_view.create_widgets()
    
    def refresh_new_image_selection(self):
        print("Refresh due to new image selection to " + str(self.current_index_var.get()))
        self.refresh_image_info()
        self.refresh_image_view()
    
    def create_control(self):
        self.button_path_selection = ctk.CTkButton(self.control_frame, text="Path Selection", command=self.changeToPathSelection)
        self.button_path_selection.place(relx=0.4, rely=0.5, anchor="e")
        # self.button_image_list = ctk.CTkButton(self.control_frame, text="Image List", command=self.changeToImageList)
        # self.button_image_list.place(relx=0.6, rely=0.5, anchor="w")
    
    def change_rating(self, rating: int):
        with ExifToolHelper() as et:
            et.set_tags(
                [self.images[self.current_index_var.get()][1], self.images[self.current_index_var.get()][1].replace(".JPG", ".CR3")],
                tags={"Rating": rating},
                params=["-P", "-overwrite_original"]
            )
        print(f"Rating of {self.images[self.current_index_var.get()][1]} changed to {rating}")
        self.refresh_image_info()
        self.refresh_image_list()
        
if __name__ == "__main__":
    from imageApp import ImageApp
    
    root = ctk.CTk()
    root.title("Image Viewer")
    app = ImageApp(root, "/mnt/daten/1 Fotos/2024/2024-04-01 Ostermontag")
    app.create_widgets()
    root.mainloop()
