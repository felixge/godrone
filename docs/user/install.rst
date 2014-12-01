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

There are currently no downloads available. You need to build from source.

Install
-------

See :ref:`installing from source <isource>`.

Uninstall
---------

The ``go get`` install process simply adds some files to your ``$GOPATH``
directory.

When you run ``godrone-util``, it stops the running copy of the factory
firmware, copies the a new executable onto the drone without touching
the factory firmware, and runs it.

If you want the original firmware back, simply disconnect and
reconnect the battery. The install is not permanent.
