# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread
import threading

i = 0

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")

lock = threading.Lock()


def function1():
	global i
	for j in range(0, 1000000):
		lock.acquire()
		i += 1
		lock.release()

def function2():
	global i
	for k in range(0, 1000000):
		lock.acquire()
		i -= 1
		lock.release()

def main():
    thread_1 = Thread(target = function1, args = (),)
    thread_2 = Thread(target = function2, args = (),)
    thread_1.start()
    thread_2.start()
    
    thread_1.join()
    thread_2.join()
    print(i)


main()
