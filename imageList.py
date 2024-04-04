import customtkinter as ctk
from tkinter import Event
from PIL import Image, ImageTk

class ImageList:
    def __init__(self, root : ctk.CTkFrame, images : list[tuple[Image.Image, str]], curr_index_var : ctk.IntVar) -> None:
        self.root = root
        self.images = images
        self.curr_index_var = curr_index_var
        
        self.create_widgets()
    
    def create_widgets(self):
        pass