import customtkinter as ctk
from CTkListbox import CTkListbox
from tkinter import Event
from PIL import Image, ImageTk
from alive_progress import alive_bar

# information in the image list
# - thumbnail
# - filename
# - rating
# - date
# - time
# - camera type
# - lens type
# - focal length
# - exposure time

class ImageList:
    def __init__(self, root : ctk.CTkFrame, images : list[tuple[Image.Image, str]], curr_index_var : ctk.IntVar, refresh_all: callable) -> None:
        self.root = root
        self.curr_index_var = curr_index_var
        self.curr_images = []
        self.refresh_all = refresh_all
        self.refresh = True
        
        # self.root.grid_rowconfigure(, weight=1)
        # self.root.grid_columnconfigure(0, weight=1) #empty
        # self.root.grid_columnconfigure(1, weight=1) #thumbnail
        # self.root.grid_columnconfigure(2, weight=3) #filename and other information
        # self.root.grid_columnconfigure(3, weight=1) #empty
        
        self.root.grid_rowconfigure(0, weight=1)
        self.root.grid_columnconfigure(0, weight=1)
        
        self.create_thumbnails(images)
        
        self.listbox = CTkListbox(self.root, command=self.on_select)
        self.listbox.grid(row=0, column=0, sticky="nsew")
        
        self.listbox.bind_all("<MouseWheel>", self.scroll)
        self.listbox.bind_all("<Button-4>", self.scroll)
        self.listbox.bind_all("<Button-5>", self.scroll)

        self.create_widgets()
        
    def create_thumbnails(self, images: list[tuple[Image.Image, str]]):
        self.images = images
        return
        with alive_bar(len(images), title="Creating thumbnails") as bar:
            for img, img_path in images:
                img = img.copy()
                img.thumbnail((128, 128), Image.LANCZOS)
                self.images.append((ImageTk.PhotoImage(img), img_path))
                bar()
    
    def create_widgets(self):
        self.listbox.delete("all")
        
        with alive_bar(len(self.images), title="Creating list") as bar:
            for i, (img, img_path) in enumerate(self.images):
                self.listbox.insert(i, img_path.split("/")[-1])
                bar()
                
        # for i, (img, img_path) in enumerate(self.images):
        #     self.listbox.insert(i, img_path.split("/")[-1])
            # ctkimage = ctk.CTkImage(img)
            # self.curr_images.append(ctkimage)
            # ctk.CTkLabel(self.root, image=ctkimage).grid(row=i, column=1, sticky="w")
            # ctk.CTkLabel(self.root, text=img_path.split("/")[-1]).grid(row=i, column=2, sticky="w")
        
        self.listbox.activate(self.curr_index_var.get())
    
    def on_select(self, selected_option):
        index = self.listbox.curselection()
        if index:
            if type(index) == tuple:
                self.curr_index_var.set(index[0])
            elif type(index) == int:
                self.curr_index_var.set(index)
            else:
                raise ValueError(f"Unknown index type: {type(index)}")#

            if self.refresh:
                self.refresh_all()
    
    def refresh_metadata(self):
        # block refreshing while setting current index because the image does not have to be reset because the refresh came from there
        self.refresh = False
        self.listbox.activate(self.curr_index_var.get())
        self.refresh = True
    
    def scroll(self, event: Event):
        if event.delta < 0 or event.num == 5:
            self.listbox._parent_canvas.yview_scroll(1, "units")
        else:
            self.listbox._parent_canvas.yview_scroll(-1, "units")