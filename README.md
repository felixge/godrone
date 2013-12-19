# godrone

An open firmware for the [Parrot AR Drone
2](http://en.wikipedia.org/wiki/Parrot_AR.Drone#AR.Drone_2.0_.282012.29)
written in Go.

## Status

This project is a work in progress, but the following things should work:

* Navboard driver
* Motorboard driver
* Complementary flight stablization filter
* Gamepad control via WebSocket / HTML5 Gamepad API
* Basic flight!

A 0.1 release will be announced soon. It will come with binary installers for
OSX, Linux and Windows.

## Roadmap

The following things still need to be implemented:

* Ultrasound height detection
* Battery status
* Reset motor emergency mode
* High level API for writing apps
* Camera Access
* Optical flow tracking for bottom camera / better hover stabilization
* Kalman filter (supposedly better performance than complementary filter)
* Parrot UDP protocol
* Parrot TCP Video Protocol
* JS scripting

## Motivation

This project is mainly a personal challenge I set for myself.

However, if it turns out well, this firmware may become a viable replacement
for the Parrot AR Drone firmware, allowing for a few interesting use cases:

* Education: Demonstrate quad copter physics, by controlling motors
  individually.
* Autonomy: Write autonomous drone software that requires no Wi-Fi connection /
  client to be connected.
* HTML client: Control the drone from the web browser of any device via HTML /
  JS / WebSockets. No need for custom apps.
* Acrobatics: Teach the drone new acrobatic tricks and allow for more aggresive
  flight.
* Hackability: Easy support for additional devices connected to the drone (GPS,
  Sensors, Lasers, etc.)
