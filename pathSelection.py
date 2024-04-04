import customtkinter as ctk

class PathSelection:
    def __init__(self, root: ctk.CTkFrame, default_path: str, folder_path_var: ctk.StringVar, on_selection: callable, scale: float = 1.5) -> None:
        self.root = root
        self.default_path = default_path
        self.folder_path_var = folder_path_var
        self.on_selection = on_selection
        
        self.scale = scale
        self.create_widgets()
    
    def create_widgets(self) -> None:
        self.main = ctk.CTkFrame(self.root, bg_color=self.root._bg_color)
        self.main.place(relx=0.5, rely=0.5, anchor=ctk.CENTER)
        # self.main_frame.grid(row=0, column=0, padx=10, pady=10)
        
        self.main_frame = ctk.CTkFrame(self.main, bg_color=self.root._bg_color)#, width=500)
        padding = 100
        # self.main_frame.grid(row=0, column=0, padx=padding, pady=padding)
        self.main_frame.pack(fill=ctk.BOTH, expand=True, padx=padding, pady=padding)
        # ctk.set_widget_scaling(self.scale)
        
        self.label_path_selection = ctk.CTkLabel(self.main_frame, text="Please select a folder with images.")
        self.label_path_selection.pack()
        
        self.button_path_selection = ctk.CTkButton(self.main_frame, text="Browse", command=self.on_select_folder)
        self.button_path_selection.pack()
        
        self.selected_folder = ctk.CTkLabel(self.main_frame, textvariable=self.folder_path_var)
        self.selected_folder.pack()
        
        self.confirm_button = ctk.CTkButton(self.main_frame, text="Confirm", command=self.on_confirm_folder)
        self.confirm_button.pack()
        
    def on_select_folder(self) -> None:
        print("Selecting folder...")
        folder = ctk.filedialog.askdirectory()
        if not folder:
            return
        self.folder_path = folder
        self.folder_path_var.set(self.folder_path)
    
    def on_confirm_folder(self) -> None:
        # ctk.set_widget_scaling(1.0) # Reset scaling before changing the state
        self.on_selection()

        