from conftest import send, recv


def test_sub_happy_path(conn):
    send(conn, "SUB sports.scores 1")
    assert recv(conn) == "+OK"


def test_sub_missing_args(conn):
    send(conn, "SUB")
    assert recv(conn).startswith("-ERR")


def test_unsub_happy_path(conn):
    send(conn, "UNSUB 1")
    assert recv(conn) == "+OK"


def test_unsub_missing_args(conn):
    send(conn, "UNSUB")
    assert recv(conn).startswith("-ERR")


def test_pub_happy_path(conn):
    conn.sendall(b"PUB sports.scores 5\r\nHello\r\n")
    assert recv(conn) == "+OK"


def test_pub_missing_args(conn):
    send(conn, "PUB")
    assert recv(conn).startswith("-ERR")


def test_ack_happy_path(conn):
    send(conn, "ACK 42")
    assert recv(conn) == "+OK"


def test_ack_missing_args(conn):
    send(conn, "ACK")
    assert recv(conn).startswith("-ERR")


def test_unknown_command(conn):
    send(conn, "BLAH foo 1")
    assert recv(conn).startswith("-ERR")
