from PIL import Image
import os

def rotate_image(img: Image) -> Image:
    rotation = img.getexif().get(0x0112)        
    if rotation == 1:
        pass
    elif rotation == 2:
        img = img.transpose(Image.FLIP_LEFT_RIGHT)
    elif rotation == 3:
        img = img.rotate(180, expand=True)
    elif rotation == 4:
        img = img.transpose(Image.FLIP_TOP_BOTTOM)
    elif rotation == 5:
        img = img.rotate(90, expand=True)
        img = img.transpose(Image.FLIP_LEFT_RIGHT)
    elif rotation == 6:
        img = img.rotate(270, expand=True)
    elif rotation == 7:
        img = img.rotate(270, expand=True)
        img = img.transpose(Image.FLIP_TOP_BOTTOM)
    elif rotation == 8:
        img = img.rotate(90, expand=True)
    else:
        print(f"Unknown rotation: {rotation}, not in https://exiftool.org/TagNames/EXIF.html")
    return img


def load_images_from_folder(folder: str):
    images = []
    files = os.listdir(folder)
    files.sort()
    files = list(filter(lambda x: x.lower().endswith(('.jpg', '.jpeg')), files))
    
    for filename in files:
        img_path = os.path.join(folder, filename)
        img = Image.open(img_path)
        images.append((img, img_path))
        
    print(f"Found {len(images)} images in the folder.")
    return images

def check_image_folder(folder: str):
    files = os.listdir(folder)
    files = list(filter(lambda x: x.lower().endswith(('.jpg', '.jpeg')), files))
    print(f"Found {len(files)} images in the folder.")
    return len(files) > 0
