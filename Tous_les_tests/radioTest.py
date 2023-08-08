import subprocess
import time
import os

command = 'rtl_fm -g 30 -f 90.4M -M wfm -s 180k -E deemp |  play -v 0.1 -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k'

radio = subprocess.Popen(command, shell=True)
time.sleep(5)
radio.terminate()
get_pid = os.popen("kill $(ps -e -o pid,comm | grep 'rtl' | awk '{print $1}')")
"""
out = get_pid.read()
print(out)
os.popen(f"kill {out}")
"""