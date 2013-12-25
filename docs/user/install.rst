Install
=======

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

* **OSX:** `Universal (32 and 64-bit) <https://github.com/felixge/godrone/releases/download/v0.1.0/godrone-v0.1.0-46-g9cb1ff5-darwin-386.zip>`_
* **Windows:** `Universal (32 and 64-bit) <https://github.com/felixge/godrone/releases/download/v0.1.0/godrone-v0.1.0-46-g9cb1ff5-windows-386.zip>`_
* **Linux:** `32-bit <https://github.com/felixge/godrone/releases/download/v0.1.0/godrone-v0.1.0-46-g9cb1ff5-linux-386.tar.gz>`_ | `64-bit <https://github.com/felixge/godrone/releases/download/v0.1.0/godrone-v0.1.0-46-g9cb1ff5-linux-amd64.tar.gz>`_

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
