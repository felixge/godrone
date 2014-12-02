First flight
============

You'll need two windows, running two commands. In one window, start the web server
for the user interface: ::

  $ godrone-ui
  2014/12/01 21:39:29.057009 Listening on: :8080

In the other window, use ``godrone-util`` to cross-compile the drone software
and send it to the drone: ::

  $ godrone-util run

Now, you are ready to fly!

1. Make sure that the drone is on a level surface.

2. Open http://127.0.0.1:8080/ in your web browser, you should see the GoDrone
   web user interface. The user interface is currently just a set of graphs
   showing the input (black) and status (red) of the drone.

3. Use the following keys to control the drone:

 - ESC: Emergency! Stop all motors. Reload the page to fly again.
 - Arrows up/down: Altitude up/down
 - Arrows left/right: Yaw
 - w/s: Pitch
 - a/d: Roll
 - p: Pause the graph
 - x: Cycle through the graphs
 - c: Calibrate (drone must be on a level surface)

Next steps
----------

First of all, congratulations, you've earned the dubious badge of honor for
having tried out version |release| of a drone software, made by amateur
robotics enthusiasts, that you just found on the interwebs! You're a truly
adventurous person.

Take all the time you need to recover from this experience, but if you feel
hungry for more, please start to get involved with the :ref:`community <community>`. Report
your experience, share your problems or demand your money back.

After that, feel free to poke around the :ref:`source <source>`, and maybe even
send a patch for that thing that annoyed you the most!
