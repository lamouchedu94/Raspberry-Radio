import sys

import Adafruit_SSD1306

from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont
import time, os, subprocess

import RPi.GPIO as GPIO # Import Raspberry Pi GPIO library




RST = None
DC = 23
SPI_PORT = 0 
SPI_DEVICE = 0

# 128x64 display
disp = Adafruit_SSD1306.SSD1306_128_32(rst=RST)


disp.begin()
disp.clear()
disp.display()

width = disp.width
height = disp.height
image = Image.new('1', (width,height))
draw = ImageDraw.Draw(image)
draw.rectangle((0,0,width,height), outline=0, fill=0)

padding = -2 
top = padding
bottom = height - padding
x = 0 
font = ImageFont.load_default()
i=0 
p = 1

def screenOff() :

    draw.rectangle((0,0,width,height), outline=0, fill = 255 )
    disp.image(image)
    disp.display()
    time.sleep(1)
    draw.rectangle((0,0,width,height), outline=0, fill = 0 )
    disp.image(image)
    disp.display()

for line in sys.stdin:
    draw.rectangle((0,0,width,height), outline=0, fill = 0 )    
    draw.text((x,top), "Fréquence :", font = font, fill = 100)
    
    draw.text((x,top+8), str(line.strip()) + " Khz", font = font, fill = 100)
    disp.image(image)
    disp.display()
    print(line)
    sys.stdout.flush()

screenOff()
