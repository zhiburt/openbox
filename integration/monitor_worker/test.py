import subprocess
import sys
import time
import socket
import requests
import json
from os.path import dirname, abspath

def projectDir():
    return dirname(dirname(dirname(abspath(__file__))))

def start_monitor():
    proc = subprocess.Popen(["make", "start_monitor"], cwd=projectDir(), stdout=subprocess.PIPE)

    return proc

def start_worker(workername):
    proc = subprocess.Popen(["make","NAME="+workername, "start_worker"], cwd=projectDir(), stdout=subprocess.PIPE)
    
    return proc 


def start_services():
    proc = subprocess.Popen(["make", "start_services"], cwd=projectDir())
    proc.wait()
    
    return proc.returncode

def stop_services():
    proc = subprocess.Popen(["make", "stop_services"], cwd=projectDir())
    proc.wait()
    
    return proc.returncode

def curl(jsn, owner):
    host_name = socket.gethostname() 
    url = "http://"+ host_name + ":8082/files"

    if jsn == None:
        url = "http://" + host_name + ":8082/files/owner/"+owner
        a = requests.get(url)
        return a.json()

    a = requests.post(url, data=json.dumps(jsn))
    print(a.json(), jsn)
    return a.json()

if __name__ == "__main__":
    cases = [
        {"json": {"name":"m","owner_id":"m"}, "is_create": True, "expected": "id"},
        {"json": {}, "is_create": True, "expected": {'error': 'request error, try to check params'}},
    ]
    
    try:
        start_services()
        time.sleep(10)
        monitor = start_monitor()
        time.sleep(10)
        worker = start_worker("worker_21332")
        time.sleep(10)
        
        for indx, case in enumerate(cases):
            print("ITERATION")
            res = []
            if case["is_create"] == True:
                res = curl(case["json"], "")
            else:
                res = curl("", case["owner"])

            if res == case["expected"]:
                print('ok - ',indx)
            elif case["expected"] == "id" and "id" in res:
                print('ok - ',indx)
            else:
                print('INVALID ',indx)
    finally:
        stop_services()
        monitor.kill()
        worker.kill()
