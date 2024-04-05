import customtkinter as ctk

class PathSelection:
    def __init__(self, root: ctk.CTkFrame, default_path: str, folder_path_var: ctk.StringVar, on_selection: callable) -> None:
        self.root = root
        self.default_path = default_path
        self.folder_path_var = folder_path_var
        self.on_selection = on_selection
        
        self.root.grid_rowconfigure(0, weight=1)
        self.root.grid_columnconfigure(0, weight=1)
        
        self.create_widgets()
    
    def create_widgets(self) -> None:
        self.main_frame = ctk.CTkFrame(self.root, bg_color=self.root._bg_color)
        self.main_frame.grid(row=0, column=0, sticky="nsew", padx=20, pady=20)
        
        self.main_frame.grid_rowconfigure((0, 1, 2, 3, 4, 5), weight=1)
        self.main_frame.grid_columnconfigure((0, 1, 2), weight=1)
        
        self.label_path_selection = ctk.CTkLabel(self.main_frame, text="Please select a folder with images.")
        self.label_path_selection.grid(row=1, column=1)
        
        self.button_path_selection = ctk.CTkButton(self.main_frame, text="Browse", command=self.on_select_folder)
        self.button_path_selection.grid(row=2, column=1)
        
        self.selected_folder = ctk.CTkLabel(self.main_frame, textvariable=self.folder_path_var)
        self.selected_folder.grid(row=3, column=1)
        
        self.confirm_button = ctk.CTkButton(self.main_frame, text="Confirm", command=self.on_confirm_folder)
        self.confirm_button.grid(row=4, column=1)
        
    def on_select_folder(self) -> None:
        print("Selecting folder...")
        folder = ctk.filedialog.askdirectory(initialdir=self.folder_path_var.get())
        if not folder:
            return
        self.folder_path = folder
        self.folder_path_var.set(self.folder_path)
    
    def on_confirm_folder(self) -> None:
        self.on_selection()

        