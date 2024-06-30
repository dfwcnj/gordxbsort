# gordxbsort



- [License](#license)
`gordxbsort` is distributed under the terms of the [MIT](https://spdx.org/licenses/MIT.html) license.


## gordxbsort

This command sorts n files using radix sort and merges these sorted
files to one output.

This is very much a work in progress. Although it passes all of the
tests that I have written, there are glaring omissions. For example,
testing for fixed length records is completely absent and may not work
at all. It has been tested for sorting 10 files but not for 1.

Witn the LANG=en_US.UTF-8, as is the default on my machine, the command
is around 2/3 as fast as BSD sort but with LANG unset, this command is
arount 1/10 as fast as BSD sort. I haven't investigated the reason, but
I suspect the use of mmap and pthreads.

I will try over time to improve the test coverage and performance.

