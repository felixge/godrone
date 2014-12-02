.. _isource:

Install from source
===================

OSX & Linux
-----------

Before getting started:

* Installed the latest version of `Go <https://golang.org/doc/install>`_
* Configure a ``$GOPATH`` (e.g. via adding ``export GOPATH="${HOME}/go"`` to
  your ``~/.profile``)
* Add ``$GOPATH/bin`` to your ``$PATH`` (e.g. via adding ``export
  PATH="${GOPATH}/bin:${PATH}"`` to your ``~/.profile``)

Run the commands below to make sure you're in good shape::

    $ go version
    go version go1.2 darwin/amd64
    $ go env GOPATH
    /Users/felix/code/go
    $ echo "${PATH}"  | grep -q "$(go env GOPATH)/bin" && echo "good" || echo "bad"
    good

Next, you will need to setup Go for cross compiling to Linux/ARM. Luckily this
only takes a minute. The fastest way is this::

    $ cd $(go env GOROOT)/src
    $ GOOS=linux GOARCH=arm ./make.bash --no-clean

If this doesn't work for you for some reason, you may try to follow Dave
Cheney's `guide for cross compiling Go
<http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1>`_
instead.

Download, compile, and install GoDrone into ``$GOPATH/bin``: ::

    $ go get github.com/felixge/godrone/cmd/godrone-util
    $ go get github.com/felixge/godrone/cmd/godrone-ui

Windows
-------

Building on windows is probably also doable, but has not been attempted yet.
Please contribute to the docs if you get it working.
