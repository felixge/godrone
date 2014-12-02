Changelog
=========

Development
-----------

The current software is under development and does not have a version number.

0.1.0
-----

**2013-12-25:** This is the first release. It is aimed at adventurous software developers, and
contains the following features:

* HTML client for controlling the drone via Gamepad. Uses WebSockets, jQuery,
  and React. Data is simply rendered into tables. There is no design yet.
* Basic navboard and motorboard drivers.
* Basic binary installer for OSX, Windows and Linux.
* Initial documentation.
* Complementary attitude filter. A kalman filter will likely provide better
  results in the future.
* Basic PID algorithm for control. This will require more tuning.
* Configuration via simple toml config file. Contributed by `gwoo
  <https://github.com/gwoo>`_.
* Logging to stdout and log file on the drone.
* Manual altitude control. The ultra sound sensor is already working, but the
  initial results with it were not good enough yet. This will require some
  filtering.
* Automatic flat trim calibration on startup. These values will need to be
  saved in the future, and the yaw will have to be reset before takeoff to
  allow moving the drone around while on the ground.
* ControlTimeout to shut off the motors if the the network connection gets
  interrupted. Prevents the drone from flying into outer space.
