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

GPIO.setwarnings(False) # Ignore warning for now
GPIO.setmode(GPIO.BCM) # Use physical pin numbering
GPIO.setup(14, GPIO.IN, pull_up_down=GPIO.PUD_DOWN) # Set pin 18 to be an input pin and set initial value to be pulled low (off)
GPIO.setup(15, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
GPIO.setup(18, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)



freq = 90400

volume = 0.02
draw.text((x,top), "frequence :", font = font, fill = 100)
draw.text((x,top+8),str(freq)+" KHz", font=font, fill=255)

command = f"rtl_fm -g 30 -f {freq}K -M wfm -s 180k -E deemp |  play -v {volume} -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
radio = subprocess.Popen(command, shell=True)

def change_freq(freq, volume) :
    os.popen("kill $(ps -e -o pid,comm | grep 'rtl' | awk '{print $1}')")
    time.sleep(1.5)
    command = f"rtl_fm -g 30 -f {freq}K -M wfm -s 180k -E deemp |  play -v {volume} -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
    radio = subprocess.Popen(command, shell=True)

def afficher_freq(freq):
    draw.rectangle((0,0,width, height), outline=0,fill=0)
    draw.text((x,top), "frequence :", font = font, fill = 255)
    draw.text((x,top+8), str(freq)+" KHz", font = font, fill = 255)
    disp.image(image)
    disp.display()


start = time.time()
tick = 1
while True :
    
    if GPIO.input(14) == GPIO.HIGH:
        freq += 100
        time.sleep(.2)
        while GPIO.input(14) == GPIO.HIGH:
            freq += 100
            time.sleep(.04)
            afficher_freq(freq)

        change_freq(freq, volume)
        afficher_freq(freq)
    
    if GPIO.input(18) == GPIO.HIGH:
        freq -= 100
        time.sleep(.2)
        while GPIO.input(18) == GPIO.HIGH:
            freq -= 100
            time.sleep(.04)
            afficher_freq(freq)
        
        change_freq(freq, volume)
        afficher_freq(freq)
    
    if GPIO.input(15) == GPIO.HIGH :
        get_pid = os.popen("kill $(ps -e -o pid,comm | grep 'rtl' | awk '{print $1}')")
        draw.rectangle((0,0,width,height), outline=0, fill=255)
        disp.image(image)
        disp.display()
        time.sleep(2)
        draw.rectangle((0,0,width, height), outline=0,fill=0)
        disp.image(image)
        disp.display()
        break
    
    if time.time() - start > 1 :
        tick *= -1 
        start = time.time()
        if tick > 0 :
            draw.rectangle((width-3,height-3,width, height), outline=0,fill=255)
        else :
            draw.rectangle((width-3,height-3,width, height), outline=0,fill=0)
        disp.image(image)
        disp.display()
    disp.image(image)
    disp.display()



"""
rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 20 -F 9 -f 90.4M |\
  redsea --feed-through |\
  play -t .s16 -r 171k -v 0.1 -c 1 -
  """