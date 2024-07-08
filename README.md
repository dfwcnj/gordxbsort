# gordxbsort



- [License](#license)
`gordxbsort` is distributed under the terms of the [MIT](https://spdx.org/licenses/MIT.html) license.


## gordxbsort

This command sorts n files using radix sort and merges these sorted
files to one output.

This is very much a work in progress. Although it passes all of the
tests that I have written, there are omissions. All of the source files
must including the test files must be inspected and cleaned up,
checking all of the errors and removing unnecessary code

Witn the LANG=en_US.UTF-8, as is the default on my machine, the command
is around 2/3 as fast as BSD sort but with LANG unset, this command is
arount 1/10 as fast as BSD sort. I haven't investigated the reason, but
I suspect the use of mmap and pthreads.

I tried making the sort phase concurrent but it only seems to make
matters worse. The insemit file useѕ insertion sort to order the merge -
it seems to be a loser.  pqchan runs concurrently but doesn't seem to
run any faster than pqread

i added an iomem argument that is used for input and removed lpo(lines
per operation)

I will try over time to improve the test coverage and performance.


