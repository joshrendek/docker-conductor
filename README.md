![](https://github.com/joshrendek/docker-conductor/blob/master/logo.png)
_Docker logo belongs to Docker Inc_

### TODO

* [ ] Health checks
* [ ] TLS connections
* [ ] Search/List containers for projects/hosts

# docker-conductor
A way to automate and orchestrate docker deployments

# Installation

`go install github.com/joshrendek/docker-conductor`

# Usage

Run `docker-conductor` inside a directory with a `conductor.yml` in it.

Flags:

```
-n, --name="": Only run the instruction with this name
-p, --project="": Only run the instruction that are apart of this project
```

#### Only deploy `test_project`

`docker-conductor -p test_project`

#### Only deploy the instruction named foobar

`docker-conductor -n foobar`

# Example conductor.yml

If you want to use a library image, just specify library infront (check example below)

``` yaml

- name: test redis
  hosts:
    - tcp://docker1.example.com:2375
  container:
    name: test-redis
    image: library/redis

- name: Descriptive Service Name
  project: test_project
  hosts:
    - tcp://docker1.example.com:2375
  container:
    name: running-container-name
    image: private.registry.example.com/yourname/your_image
    environment:
      - FOOBAR=baz
    ports:
      80/tcp: 8080
    volumes:
      - /tmp:/tmp
    dns:
      - 8.8.8.8

- name: Descriptive Service Name 2
  hosts:
    - tcp://docker1.example.com:2375
  container:
    name: foobar-baz
    image: private.registry.example.com/yourname/foobar_baz_image
```

# License
```
The MIT License (MIT)

Copyright (c) 2015 Josh Rendek

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
