[![codecov](https://codecov.io/gh/zhiburt/openbox/branch/master/graph/badge.svg)](https://codecov.io/gh/zhiburt/openbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhiburt/openbox)](https://goreportcard.com/report/github.com/zhiburt/openbox)
[![Build Status](https://travis-ci.org/zhiburt/openbox.svg?branch=master)](https://travis-ci.org/zhiburt/openbox)
[![Godoc](https://godoc.org/github.com/zhiburt/openbox?status.svg)](https://godoc.org/github.com/zhiburt/openbox)
# Openbox

Openbox this is a high availeble file service.

<p align="center">
  <img src="../assets/images/openbox.png?raw=true">
</p>

structure and alghorithm of this one you can see below

### High architect

* The monitor is the one which the hiest on the picture
* The worker is the server that will store files, that are all others nodes

#### Algorithm

1. user send file to monitor
2. monitor send one to queue
3. at the time workers watching queues
4. if user need file, request goes directly to worker's queue

[![Structure](../assets/images/structure_openbox.png?raw=true)]

### Usage

as for usage, there's no any UI parts of this project:(

one attempt have been done but, that's not thing which be worth to push to master
