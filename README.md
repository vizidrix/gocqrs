[![GoDoc](https://godoc.org/github.com/vizidrix/gocqrs?status.png)](https://godoc.org/github.com/vizidrix/gocqrs)
[![Build Status](https://drone.io/github.com/vizidrix/gocqrs/status.png)](https://drone.io/github.com/vizidrix/gocqrs/latest)

gocqrs
====================

A suite of support interfaces and structs to aide in developing a CQRS infrastructure in Golang.

Header structure allows for compact storage and fast comparison of message data while also allowing a high degree of partitioning across Applications and Domains.

Application partitions enable domains to be silo'd for multi-tennancy or other split purpose.

Domain partitions maintain a clear deliniation between boundaries both logically and on disk.
