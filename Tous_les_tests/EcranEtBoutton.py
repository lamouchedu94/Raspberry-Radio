
#import Adafruit_GPIO.SPI as SPI
import Adafruit_SSD1306

from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont
import time

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

GPIO.setwarnings(False) # Ignore warning for now
GPIO.setmode(GPIO.BCM) # Use physical pin numbering
GPIO.setup(14, GPIO.IN, pull_up_down=GPIO.PUD_DOWN) # Set pin 18 to be an input pin and set initial value to be pulled low (off)
GPIO.setup(15, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
GPIO.setup(18, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)

draw.text((x,top), "Valeur :", font = font, fill = 100)
draw.text((x,top+8),"0", font=font, fill=255)
while True : 
    

    if GPIO.input(14) == GPIO.HIGH:
        i+=1
        draw.rectangle((0,0,width,height), outline=0, fill=0)
        draw.text((x,top), "Valeur :", font = font, fill = 255)
        draw.text((x,top+8), str(i), font = font, fill = 255)
        time.sleep(0.05)
        print("Button + was pushed!")
    if GPIO.input(18) == GPIO.HIGH:
        i-=1
        draw.rectangle((0,0,width,height), outline=0, fill=0)
        draw.text((x,top), "Valeur :", font = font, fill = 255)
        draw.text((x,top+8), str(i), font = font, fill = 255)
        time.sleep(0.05)
        print("Button - was pushed!")
    if GPIO.input(15) == GPIO.HIGH :
        draw.rectangle((0,0,width,height), outline=0, fill=255)
        disp.image(image)
        disp.display()
        time.sleep(2)
        draw.rectangle((0,0,width, height), outline=0,fill=0)
        disp.image(image)
        disp.display()
        break

    disp.image(image)
    disp.display()

