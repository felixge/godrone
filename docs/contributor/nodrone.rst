Testing with No Drone
=====================

Sometimes it is not convenient to have a drone buzzing aorund your head
while you are testing something.

To run ``godrone`` on your laptop, simply use ``go build`` to build it.
Run it with the ``-dummy`` argument to tell it to not connect to the
drone's navboard and motorboard. It will try to listen on port 80.
Since you probably don't want to have to run it as root, use the
``-addr=:8000`` argument. To see what would be written to the motor board,
use ``-verbose=2``.

To get the UI to connect to it, use the following URL:
        http://127.0.0.1:8080/?127.0.0.1:8000

