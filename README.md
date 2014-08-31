# h2check

h2check is a conformance test tool for HTTP/2 servers.  
This tool supports [draft-ietf-httpbis-http2-14](http://tools.ietf.org/html/draft-ietf-httpbis-http2-14).

## Install

```
$ git clone https://github.com/summerwind/h2check.git
$ cd h2check
$ npm install
```
*h2check requires [Node.js](http://nodejs.org/), [npm](https://www.npmjs.org/) and [Git](http://git-scm.com/).*

## Usage

```
$ ./bin/h2check --help

  Usage: h2check [options]

  Options:

    -h, --help          output usage information
    -p, --port <value>  target port
    -h, --host <value>  target host
``` 

## License

The MIT License

Copyright (c) 2014 Moto Ishizawa

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


