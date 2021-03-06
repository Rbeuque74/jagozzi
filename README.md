jagozzi [![Go Report Card](https://goreportcard.com/badge/github.com/rbeuque74/jagozzi)](https://goreportcard.com/report/github.com/rbeuque74/jagozzi) [![Build Status](https://travis-ci.org/rbeuque74/jagozzi.png?branch=master)](https://travis-ci.org/rbeuque74/jagozzi) [![Coverage Status](https://coveralls.io/repos/github/rbeuque74/jagozzi/badge.svg?branch=master)](https://coveralls.io/github/rbeuque74/jagozzi?branch=master) [![GitHub release](https://img.shields.io/github/release/rbeuque74/jagozzi.svg)](https://github.com/rbeuque74/jagozzi/releases)
==============================

jagozzi is a light monitoring daemon for severals service in order to report results checks to a remote NSCA server.

This program is a Golang clone of [sauna](https://github.com/NicolasLM/sauna) that will parse the same configuration file format.

Services included
-----------------

- Supervisor
- Command
- Processes
- HTTP
- Marathon

Installation
------------

jagozzi can be installed using this command:

```
go install github.com/rbeuque74/jagozzi
```

Screenshot
----------

When jagozzi is sstarted with `-display` option, a "display" consumer is created and render status of configurated checks.

[![Display mode](https://image.ibb.co/jmeS6T/rsz_capture_d_cran_de_2018_05_20_00_18_00.png)](https://image.ibb.co/jmeS6T/rsz_capture_d_cran_de_2018_05_20_00_18_00.png)


License
-------

MIT, see [LICENSE](LICENSE)
