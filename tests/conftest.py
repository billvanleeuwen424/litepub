import os
import socket
import subprocess
import time
import pytest

BROKER_BINARY = os.environ.get("LITEPUB_BINARY", "../bin/litepub")
BROKER_HOST = "127.0.0.1"
BROKER_PORT = 8080


@pytest.fixture(scope="session")
def broker():
    # Pass GOCOVERDIR explicitly so coverage data is written even if the
    # variable was set after the pytest process started (belt-and-suspenders).
    env = os.environ.copy()
    proc = subprocess.Popen([BROKER_BINARY], env=env)
    time.sleep(0.2)  # wait for broker to bind
    yield
    # SIGTERM is caught by main.go's signal handler, which flushes coverage
    # data and calls os.Exit(0), triggering Go's exit hooks.  wait() ensures
    # the flush is complete before go tool covdata runs.
    proc.terminate()
    proc.wait()


@pytest.fixture
def conn(broker):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((BROKER_HOST, BROKER_PORT))
    s.settimeout(2)
    yield s
    s.close()


def send(s, line):
    s.sendall((line + "\r\n").encode())


def recv(s):
    buf = b""
    while not buf.endswith(b"\r\n"):
        buf += s.recv(1)
    return buf.decode().strip()
