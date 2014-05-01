sieve ![Project status](http://img.shields.io/status/experimental.png?color=red)
=====

A command line utility for graphing piped data. Sieve is used when you have
data output from one program and you need to quickly and interactively inspect
it without needing the full power of gnuplot.


## Usage

To install sieve, simply run:

```sh
$ go install github.com/benbjohnson/sieve/...
```

Then run it from the command line:

```sh
$ cat mydata.json | sieve
Listening on http://localhost:6900
```

Open [http://localhost:6900](http://localhost:6900) in your browser and you'll
see the data that was piped from `mydata.json`.


## Example

To see sieve in action you can use the random walk generator bundled with sieve:

```sh
$ sievegen | sieve
```

Then open [http://localhost:6900](http://localhost:6900) in your browser can watch the random walk.
