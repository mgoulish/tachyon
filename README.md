# tachyon
Go Society of Mind experiment

Tachyon is a system that provides higher-level communication patterns implemented on top of Go channels. These patterns operate at high speed to support the interaction of many goroutines in more complex ways than simple point-to-point communication.

For example: topics (each message fans out to many consumers), and bulletin boards (message stays around so that consumers can browse at any time later).

These communication patterns are used by many interacting processes -- called Abstractors -- to implement a complex machine vision system.  The goal is to show intelligent behavior which has not been explicitly planned by the programmer emerging from the interaction of many simple Abstractors.

The subject area on on which this system is being developed is sequences of images from space flight videos. The goal at this level is to create vision systems powerful enough to permit autonomous operation of spacecraft during exploration, mining, materials processing, and construction operations when they are too distant to permit real-time human control.

