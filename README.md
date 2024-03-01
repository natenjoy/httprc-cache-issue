# httprc-cache-issue

Minimal repro for httprc cache issue. This repo will be obsolete once the issue is fixed.

The issue of interest is that cache.Get is blocked while cache.Refresh happens: ([v1.0.4 code here](https://github.com/lestrrat-go/httprc/blob/v1.0.4/cache.go#L147-L159)).

If I understand correctly, we could only have the routes that can end up fetching acquire
the semaphore. Then, the worst that could happen in cache.Get while a refresh is happening
is that it gets stale data. Because of the read lock, there should be no concurrency issue.

Run example with `go run main.go`

The output will be something like:

```
SETUP: spinning up server
SETUP: registering servers with cache
SETUP: sleeping for 2 hours to watch issue in terminal
CLIENT: cache get at   Mar  1 11:55:41.489
SERVER: GET received
SERVER: will time out
CLIENT: cache get err failed to fetch "http://0.0.0.0:41234": failed to fetch "http://0.0.0.0:41234": Get "http://0.0.0.0:41234": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
CLIENT: cache value at Mar  1 11:55:47.494: nil
CLIENT: cache get at   Mar  1 11:55:48.495
CLIENT: cache value at Mar  1 11:55:48.495: nil
CLIENT: cache get at   Mar  1 11:55:49.496
CLIENT: cache value at Mar  1 11:55:49.496: nil
CLIENT: cache get at   Mar  1 11:55:50.497
CLIENT: cache value at Mar  1 11:55:50.497: nil
CLIENT: cache get at   Mar  1 11:55:51.498
CLIENT: cache value at Mar  1 11:55:51.498: nil
CLIENT: cache get at   Mar  1 11:55:52.499
CLIENT: cache value at Mar  1 11:55:52.499: nil
CLIENT: cache get at   Mar  1 11:55:53.499
SERVER: GET received
CLIENT: cache value at Mar  1 11:55:53.501: count 1
CLIENT: cache get at   Mar  1 11:55:54.502
CLIENT: cache value at Mar  1 11:55:54.502: count 1
CLIENT: cache get at   Mar  1 11:55:55.503
CLIENT: cache value at Mar  1 11:55:55.503: count 1
CLIENT: cache get at   Mar  1 11:55:56.504
CLIENT: cache value at Mar  1 11:55:56.504: count 1
CLIENT: cache get at   Mar  1 11:55:57.505
CLIENT: cache value at Mar  1 11:55:57.505: count 1
CLIENT: cache get at   Mar  1 11:55:58.506
CLIENT: cache value at Mar  1 11:55:58.506: count 1
SERVER: GET received
CLIENT: cache get at   Mar  1 11:55:59.507
CLIENT: cache value at Mar  1 11:55:59.507: count 2
CLIENT: cache get at   Mar  1 11:56:00.508
CLIENT: cache value at Mar  1 11:56:00.508: count 2
CLIENT: cache get at   Mar  1 11:56:01.509
CLIENT: cache value at Mar  1 11:56:01.509: count 2
CLIENT: cache get at   Mar  1 11:56:02.510
CLIENT: cache value at Mar  1 11:56:02.510: count 2
CLIENT: cache get at   Mar  1 11:56:03.511
CLIENT: cache value at Mar  1 11:56:03.511: count 2
CLIENT: cache get at   Mar  1 11:56:04.512
CLIENT: cache value at Mar  1 11:56:04.512: count 2
SERVER: GET received
CLIENT: cache get at   Mar  1 11:56:05.513
CLIENT: cache value at Mar  1 11:56:05.513: count 3
CLIENT: cache get at   Mar  1 11:56:06.515
CLIENT: cache value at Mar  1 11:56:06.515: count 3
CLIENT: cache get at   Mar  1 11:56:07.516
CLIENT: cache value at Mar  1 11:56:07.516: count 3
CLIENT: cache get at   Mar  1 11:56:08.517
CLIENT: cache value at Mar  1 11:56:08.517: count 3
CLIENT: cache get at   Mar  1 11:56:09.518
CLIENT: cache value at Mar  1 11:56:09.518: count 3
CLIENT: cache get at   Mar  1 11:56:10.519
CLIENT: cache value at Mar  1 11:56:10.520: count 3
SERVER: GET received
SERVER: will time out
CLIENT: cache get at   Mar  1 11:56:11.521
CLIENT: cache value at Mar  1 11:56:16.492: count 3
CLIENT: cache get at   Mar  1 11:56:17.492
CLIENT: cache value at Mar  1 11:56:17.492: count 3
CLIENT: cache get at   Mar  1 11:56:18.493
CLIENT: cache value at Mar  1 11:56:18.493: count 3
CLIENT: cache get at   Mar  1 11:56:19.494
CLIENT: cache value at Mar  1 11:56:19.494: count 3
CLIENT: cache get at   Mar  1 11:56:20.494
CLIENT: cache value at Mar  1 11:56:20.494: count 3  
```

To preview a proposed fix for the issue, clone the following into a directory that is sibling to the bug repro, e.g. `git clone -b acquire-sem-on-fetch https://github.com/natenjoy/httprc`. Ensure

Then, add the following `go.work` file to the root of this repo, and run `go run main.go`

```
go 1.22.0

use (
	.
	../natenjoy-httprc
)
```
