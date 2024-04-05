import customtkinter as ctk
from tkinter import Event
from PIL import Image, ImageTk
from util import rotate_image

SCALE_FACTOR = 1.2

class ImageView:
    def __init__(self, root : ctk.CTkFrame, images : list[tuple[Image.Image, str]], curr_index_var : ctk.IntVar) -> None:
        self.root = root
        self.images = images
        self.curr_index_var = curr_index_var
        
        self.photo_image = None
        self.zoomed_in = False
        self.scale = 1.0

        self.canvas = ctk.CTkCanvas(
            root,
            background='#3B3B3B',
            borderwidth=0,
            highlightthickness=0,
            width=1920,
            height=1080
        )
        
        self.canvas.bind("<Button-4>", self.zoom)
        self.canvas.bind("<Button-5>", self.zoom)
        self.canvas.bind("<MouseWheel>", self.zoom)
        self.canvas.bind("<ButtonPress-1>", self.start_drag)
        self.canvas.bind("<B1-Motion>", self.drag)
        
        self.canvas.bind("<Double-Button-1>", self.double_click_zoom)

        self.create_widgets()
    
    def double_click_zoom(self, event: Event):
        self.scale = 1.0
        self.zoomed_in = False
        self.create_widgets()
        return
    
    def zoom(self, event: Event):
        if event.delta > 0 or event.num == 4:
            self.scale *= SCALE_FACTOR
        else:
            self.scale /= SCALE_FACTOR
        self.create_widgets()

    def start_drag(self, event):
        self.start_x = event.x
        self.start_y = event.y

    def drag(self, event):
        delta_x = event.x - self.start_x
        delta_y = event.y - self.start_y
        self.start_x = event.x
        self.start_y = event.y
        self.canvas.move(self.image_id, delta_x, delta_y)

    def resizing(self, event: Event):
        if event.width > 1 and event.height > 1:
            self.create_widgets()
                    
    def create_widgets(self):
        self.canvas.delete(ctk.ALL)
        self.canvas.pack(fill=ctk.BOTH, expand=True)
        self.canvas.focus_set()


        self.root.update()
        height = self.root.winfo_height()
        width = self.root.winfo_width()

        img, img_path = self.images[self.curr_index_var.get()]
        img = img.copy()

        img = rotate_image(img)
        
        img.thumbnail((width * self.scale, height * self.scale), Image.LANCZOS)

        self.photo_image = ImageTk.PhotoImage(img)
        self.image_id = self.canvas.create_image(width / 2, height / 2, anchor=ctk.CENTER, image=self.photo_image)

if __name__ == "__main__":
    from imageApp import ImageApp
    
    root = ctk.CTk()
    root.title("Image Viewer")
    app = ImageApp(root, "~/Pictures")
    app.create_widgets()
    root.mainloop()
