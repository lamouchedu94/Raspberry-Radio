import sys

import Adafruit_SSD1306

from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont
import time, os, subprocess

import RPi.GPIO as GPIO # Import Raspberry Pi GPIO library
import threading



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

NAME = ""
FREQ = ""

def screenOff() :

    draw.rectangle((0,0,width,height), outline=0, fill = 255 )
    disp.image(image)
    disp.display()
    time.sleep(1)
    draw.rectangle((0,0,width,height), outline=0, fill = 0 )
    disp.image(image)
    disp.display()

def screen(dico):
    draw.rectangle((0,0,width,height), outline=0, fill = 0 ) 
    draw.text((x,top), "Fréquence :", font = font, fill = 100)
    draw.text((x,top+8), dico["freq"], font = font, fill = 100)
    draw.text((x,top+16), dico["name"], font = font, fill = 100)
    disp.image(image)
    disp.display()

        

dico_info = {"freq":"", "name":""}

thread_screen = threading.Thread(target=screen, args=(dico_info))
#thread_screen.start()
for line in sys.stdin:
    
    if line[0:4] == "freq" : 
        dico_info["freq"] = line[5:11] + " Khz"
    
    if line[0:4] == "name" : 
        if len(line) >= 6 :
            car = line[len(line)-2]
            dico_info["name"] = line[5:]
            
    screen(dico_info)
    print(dico_info)
    sys.stdout.flush()
    



screenOff()
