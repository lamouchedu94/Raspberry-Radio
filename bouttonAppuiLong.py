import RPi.GPIO as GPIO 
import time
GPIO.setwarnings(False) # Ignore warning for now
GPIO.setmode(GPIO.BCM) # Use physical pin numbering
GPIO.setup(14, GPIO.IN, pull_up_down=GPIO.PUD_DOWN) # Set pin 18 to be an input pin and set initial value to be pulled low (off)
GPIO.setup(15, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
GPIO.setup(18, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
val = 0
start = time.time()
while True: # Run forever
    if GPIO.input(18) == GPIO.HIGH:
        #première appui
        val += 1
        print("Button was pushed!")
        time.sleep(.1)
        while GPIO.input(18) == GPIO.HIGH:
            val += 1
            time.sleep(.1)
        print(val)
    
        