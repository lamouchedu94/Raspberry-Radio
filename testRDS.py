"""
import subprocess
import time
cmd1 = "rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 20 -F 9 -f 96M" 
proc1 = subprocess.Popen(cmd1, stdout=subprocess.PIPE, shell=True)
stdout1,_ = proc1.communicate()

cmd2 = "play -v 0.1 -r 180k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k"
proc2 = subprocess.Popen(cmd2, stdin=subprocess.PIPE, stdout=subprocess.PIPE, shell=True)
stdout2, _ = proc2.communicate(input=stdout1)
"""

import subprocess
import threading
import time

myPopen = subprocess.Popen("rtl_fm -M fm -l 0 -A std -p 0 -s 171k -g 30 -F 9 -f 90.4M" , shell = True, stdout = subprocess.PIPE)
cmd2 = subprocess.Popen("play -v 0.03 -r 171k -t raw -e s -b 16 -c 1 -V1 - lowpass 16k", stdin=myPopen.stdout , shell = True)
#cmd3 = subprocess.Popen("redsea --show-partial | jq '.partial_ps'",stdin=myPopen.stdout , shell = True,)
i = 0 
while True :
    i += 1
    print(i)

