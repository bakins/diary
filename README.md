Diary is a simple go logger.

[![GoDoc](https://godoc.org/github.com/bakins/diary?status.svg)](https://godoc.org/github.com/bakins/diary) [![Build Status](https://travis-ci.org/bakins/diary.svg?branch=master)](https://travis-ci.org/bakins/diary)

Consider this beta quality.

Diary is an opinionated logger. It:

- logs to Stdout by default. Though any io.writer can be used.
- logs in json format **only**
- does not handle concurrency or locking. This is unneeded for `os.Stdout` and `os.Stderr`.  The Typical use case is to create a logger and then it is read only, except for actually logging.  If you need locking, implement it in a wrapping structure and/or in your writer.

The primary goal was to get the developer interface correct. Lots of loggers out there have horrible developer experience, IMO.  The innnards of this can be optimized I'm sure.



