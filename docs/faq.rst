Frequently Asked Questions
==========================

I need help, where can I get support?
-------------------------------------

The best way to get help is by posting to the :ref:`godrone mailing list
<godrone-list>` or joining our :ref:`IRC Channel <godrone-irc>`.

Please keep in mind that this is a volunteer effort, so try to be nice.

I want to help, how can I get involved?
---------------------------------------

If you're a designer or a software developer who knows HTML, CSS, JavaScript or
Go, the best way to help is by sending pull requests on :ref:`GitHub <source>`
. There are many :ref:`open issues <tracker>` and :ref:`ideas <goals>`, but any
proposed changes will be considered.

That being said, it's always a good idea to post to the :ref:`godrone-dev
mailing <godrone-dev-list>` to discuss things before getting started.

Additionally help with the documentation and answering questions on the mailing
list would be really appreciated.


Where can I report bugs or propose features?
--------------------------------------------

If you're not sure that what you found is a bug, and you'd like to discuss the
problem first, please post to the :ref:`godrone mailing list <godrone-list>`.
Otherwise feel free to directly post on our :ref:`issue tracker <tracker>`.

Features are usually best discussed on the mailing list and added to the
tracker once there is consensus.

Does GoDrone provide any features not provided by the official firmware?
------------------------------------------------------------------------

Not a lot at this point, but being able to control the drone with a gamepad
using your web browser is nice. It also seems that the official firmware puts
some limits on speed, especially vertically, so you may find that GoDrone will
help you with crashing your drone in more spectacular ways : ).

.. _goals:

What are the goals of the project?
----------------------------------

The project was started as a personal challenge to write a firmware for the
Parrot AR Drone in Go. Now that this is achieved, there are many possibilities,
and it will really depend on whatever users and contributors are most
interested in. Here are a few ideas:

* Improved stabilization and control algorithms (kalman filter, optical flow,
  etc.)
* Video drivers and support for streaming video to the HTML client
* Compatibility with the various `mobile apps
  <http://ardrone2.parrot.com/apps/>`_ and `NodeCopter
  <http://nodecopter.com/>`_.
* Improved HTML client, e.g. usability, design, graphs, WebGL, support for
  mobile devices
* GPS support, including the ability to plan missions in the HTML client on a
  map
* 3G support, to allow controlling the drone over cellular networks
* Education applications. e.g. a specialized HTML client to demonstrate
  quadcopter physics
* Support for uploading JavaScript scripts that run on the drone and allow
  users to create simple applications.
* Mounting the AR Drone electronics on custom frames, attaching custom sensors,
  and maybe even using different motors.
* Encourage vendors such as Parrot to embrace a full open source strategy.

Why is the firmware written in Go?
----------------------------------

Originally to see if a firmware like this could be written in a high level
garbage collected language.

Go seemed like a good choice because it has great support for cross-compiling,
concurrency, systems programming and is just a very pleasant language to work
with.

Isn't Go unsuitable for real-time applications like this?
---------------------------------------------------------

This question is one of the reasons this project exists. Go uses a
stop-the-world garbage collector that does not provide any real-time guarantees
[#gc]_, so it's certainly not a perfect fit for a flying robot.

However, for all practical purposes the GC just needs to keep up with the
stabilization loop which runs at 200 Hz. This means that GC pauses below 5ms
have no impact on performance. Longer pauses will degrade stabilization
performance, but the tolerance threshold may be up to a second depending on
altitude and the situation.

Considering that stabilization cannot be guaranteed due to environmental
factors to begin with, it will be interesting to see if drone vendors will make
similar compromises for reducing the costs of software development , or if
governments will provide detailed software architecture regulations for
commercial drones.

Given that the AR Drone is a very light weight toy that has an extremely low
chance of causing direct harm, the GoDrone project will continue to use the
current approach for now. However, if problems are observed, or the project
becomes more popular than expected, the plan is to rewrite the stabilization
loop in C, run it on a separate thread with strong scheduling guarantees, and
use some form of IPC to communicate with it.

.. _source:

Where can I get the source code?
--------------------------------

The source code is available on GitHub. https://github.com/felixge/godrone

What is the license is GoDrone release under?
---------------------------------------------

GoDrone is licensed under the `AGPLv3 license
<https://github.com/felixge/godrone/blob/master/LICENSE.txt>`_.

This basically means that any derived software products will have to be
licensed under the same license, and that their source code needs to be made
available.

The license was chosen to ensure that the GoDrone will always remain free
software. Contributors are not asked to sign a CLA, so there will be no dual
licensing model in the future.

.. [#gc] `What kind of Garbage Collection does Go use? <http://stackoverflow.com/questions/7823725/what-kind-of-garbage-collection-does-go-use/7824353#7824353>`_
