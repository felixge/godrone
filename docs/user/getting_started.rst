Getting Started
===============

Disclaimer
----------

Please be careful. Ideally you should experiment with this firmware inside a
large indoor space. If you go outside, please try to avoid proximity to streets
and people. The worst case scenario is your drone causing a traffic accident,
and you'll have nobody to blame but yourself.

Installing this software may also void your warranty and cause damage to your
drone. But don't worry, you get back the original firmware by reconnecting the
battery, and it's very easy to fix your drone if you crash it.

Download
--------

The latest version is |release|. You can download it here:

* **OSX:** `Universal (32 and 64-bit) <https://github.com/felixge/godrone/releases/download/v0.0.1/godrone-v0.0.1-46-g9cb1ff5-darwin-386.zip>`_
* **Windows:** `Universal (32 and 64-bit) <https://github.com/felixge/godrone/releases/download/v0.0.1/godrone-v0.0.1-46-g9cb1ff5-windows-386.zip>`_
* **Linux:** `32-bit <https://github.com/felixge/godrone/releases/download/v0.0.1/godrone-v0.0.1-46-g9cb1ff5-linux-386.tar.gz>`_ | `64-bit <https://github.com/felixge/godrone/releases/download/v0.0.1/godrone-v0.0.1-46-g9cb1ff5-linux-amd64.tar.gz>`_

Install
-------

1. Download one of the archives from above.
2. Extract the archive you downloaded.
3. Place your drone on a level surface, ready for takeoff, and connect the
   battery.
4. Connect your computer to your drone's WiFi network.
5. Double-click the ``deploy`` binary inside the folder extracted from the
   archive.

That's it! The ``deploy`` binary will upload GoDrone via FTP and then telnet
into your drone to start it.

Uninstall
---------

If you want the original firmware back, simply reconnect the battery. The
install is currently not permanent. This may change if the firmware becomes
more mature in the future.

First flight
------------

1. Make sure that the drone was on a level surface when you executed the
   ``deploy`` program, and that you have not rotated it since. If not, please
   run the ``deploy`` program again. Otherwise the drone will try to correct
   for the rotation on takeoff, which may lead to a crash right away.
2. Open http://192.168.1.1/ in your web browser, you should see the GoDrone
   web user interface. But don't expect too much, it's not pretty yet.
3. Connect a gamepad to your computer. Right now the only tested model is this
   `PS3-style controller
   <http://www.amazon.de/Gamepad-Vibration-Controller-schwarz-Windows/dp/B00BUNOOHQ/ref=sr_1_1?ie=UTF8&qid=1387067500&sr=8-1&keywords=mac+gamepad>`_
   (please share your results of using other controllers).
4. Press the B0 button (X) to turn on the engines.
5. Move the right joystick up / down to control thrust.
6. Use the left joystick to control the roll and pitch angle of your drone,
   allowing you to fly in any desired direction.
7. Press the B0 button (X) again to turn off the engines.

Next steps
----------

First of all, congratulations, you've earned the dubious badge of honor for
having tried out version |release| of a drone software, made by amateur
robotics enthusiasts, that you just found on the interwebs! You're a truly
adventurous person.

Take all the time you need to recover from this experience, but if you feel
hungry for more, please start to get involved with the :doc:`community`. Report
your experience, share your problems or demand your money back.

After that, feel free to poke around the :ref:`source <source>`, and maybe even
send a patch for that thing that annoyed you the most!
