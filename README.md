# pipoint - a PX4 based camera pointer.

Some time ago I set up a camera, pointed it at the sky, and recorded
as I flew my model plane about.  It was quite cool, but the plane
covers so much area that most of the video was of blue sky.

PiPoint solves this problem by automatically pointing a ground based
camera at the rover using GPS, a telemtry link, and pan/tilt unit.

## Build

* See `ansible/` for rules to set up the Raspberry Pi.
* See `etc/` for files used on the PX4 or Raspberry Pi.
* See `Makefile` for shortcuts to build pipoint itself.

More instructions to come!

# Note
This is not an official Google product.

-- Michael Hope <michaelh@juju.net.nz> <mlhx@google.com>
