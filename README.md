CQRS (Command Query Responsibility Separation) / ES (Event Sourcing) utilities for Go

[![Build Status](https://drone.io/github.com/Vizidrix/gocqrs/status.png)](https://drone.io/github.com/Vizidrix/gocqrs/latest)

## Introduction ##

----

> CQRS is not a framework, nor does it need one
>> I built one anyways...

  - Gotta have [Go]: http://golang.org

----

This library is intended to provide some reusable constructs when building CQRS systems

## Features ##

Aggregate
Command
Event


A few unit tests working against a sample implementation.  So far the example focuses on Aggregate loading but will soon be expanded into Event Sourcing and Pub/Sub.

## Getting Started ##

1\. Add the correct import for your project.

```go
import (
	"github.com/vizidrix/gocqrs"
)
```

2\. Start using.

See examples direction and gocqrs_test.go for more detailed information.  (Real docs pending)

# Goals #
- Provide a reusable abstraction over common CQRS concepts
- Enable quick bootstrapping of new domains
- Hone in on an optimum api for interacting with CQRS/ES in Go
- Enable isolated benchmarking of core library to find new ways of organizing the interactions

----

Version
----
0.1.0 ish

Tech
----

* [Go] - Golang.org
* [GOCQRS] - CQRS (Command Query Responsibility Separation) / ES (Event Sourcing) utilities for Go

License
----

https://github.com/Vizidrix/gocqrs/blob/master/LICENSE

----
## Edited
* 11-July-2013	initial release

----
## Credits
* Vizidrix <https://github.com/organizations/Vizidrix>
* Perry Birch <https://github.com/PerryBirch>