import customtkinter as ctk
from tkinter import Event
from util import load_images_from_folder
from enum import Enum
from pathSelection import PathSelection
from imageDisplay import ImageDisplay

class AppState(Enum):
    PATH_SELECTION = 1
    IMAGE_LIST = 2
    IMAGE_VIEW = 3

class ImageApp:
    def __init__(self, root: ctk.CTk, default_folder: str) -> None:
        ctk.set_default_color_theme("dark-blue")
        ctk.set_appearance_mode("system")
        
        self.root = root
        self.root.minsize(500, 200)
        self.root.title("Image Viewer")
        self.state = AppState.PATH_SELECTION
        
        self.root_frame = ctk.CTkFrame(self.root)
        self.root_frame.pack(fill=ctk.BOTH, expand=True)
        
        self.folder_name = None
        self.folder_name_var = ctk.StringVar()
        self.folder_name_var.set(default_folder)
        
        self.root.bind("<Escape>", self.on_escape)
        self.root.bind("<F11>", self.on_fullscreen)
        self.root.bind("<Configure>", self.on_resize)
        
        self.root.bind("<Left>", lambda event: self.previous_image())
        self.root.bind("<Up>", lambda event: self.previous_image())
        self.root.bind("<Right>", lambda event: self.next_image())
        self.root.bind("<Down>", lambda event: self.next_image())
        
        self.root.bind("1", lambda event: self.change_rating(1))
        self.root.bind("2", lambda event: self.change_rating(2))
        self.root.bind("3", lambda event: self.change_rating(3))
        self.root.bind("4", lambda event: self.change_rating(4))
        self.root.bind("5", lambda event: self.change_rating(5))
        self.root.bind("0", lambda event: self.change_rating(0))
    
    def on_escape(self, event: Event) -> None:
        self.root.quit()
    
    def on_fullscreen(self, event: Event) -> None:
        self.root.attributes("-fullscreen", not self.root.attributes("-fullscreen"))
    
    def on_resize(self, event: Event) -> None:
        # self.root_frame.configure(height=self.root.winfo_height())
        # self.create_widgets() # if this is active, python crashes with a segmentation fault
        pass
    
    def next_image(self) -> None:
        if self.state == AppState.IMAGE_VIEW:
            self.image_view_parent.next_image()
    
    def previous_image(self) -> None:
        if self.state == AppState.IMAGE_VIEW:
            self.image_view_parent.previous_image()
    
    def change_rating(self, rating: int) -> None:
        if self.state == AppState.IMAGE_VIEW:
            self.image_view_parent.change_rating(rating)
    
    
    def clean_root_frame(self) -> None:
        children = list(self.root_frame.children.keys()).copy()
        for child_id in children:
            self.root_frame.children[child_id].destroy()
    
    def create_widgets(self) -> None:
        self.clean_root_frame()
                
        if self.state == AppState.PATH_SELECTION:
            self.create_path_selection()
        elif self.state == AppState.IMAGE_LIST:
            self.create_image_list()
        elif self.state == AppState.IMAGE_VIEW:
            self.create_image_view()
        else:
            raise ValueError(f"Unknown state: {self.state}")
    
    def create_path_selection(self) -> None:
        self.clean_root_frame()
        self.root.title("Image Viewer - Select Folder")
        self.path_selection = PathSelection(self.root_frame, self.folder_name, self.folder_name_var, self.on_select_folder)

    def on_select_folder(self) -> None:
        self.root.attributes('-zoomed', False)
        
        self.folder_name = self.folder_name_var.get()
        print(f"Selected folder {self.folder_name}")
        self.images = load_images_from_folder(self.folder_name)
        self.state = AppState.IMAGE_VIEW
        self.create_widgets()
    
    def create_image_list(self) -> None:
        self.clean_root_frame()
        self.root.title("Image Viewer - Image List")
        self.root.attributes('-zoomed', True)
        pass
    
    def create_image_view(self) -> None:
        self.clean_root_frame()
        self.root.title("Image Viewer - Image View")
        self.root.attributes('-zoomed', True)
        self.image_view_parent = ImageDisplay(self.root_frame, self.images, changeToPathSelection=self.create_path_selection, changeToImageList=self.create_image_list)
        

if __name__ == "__main__":
    root = ctk.CTk()
    root.title("Image Viewer")
    app = ImageApp(root, "~/Pictures")
    app.create_widgets()
    root.mainloop()