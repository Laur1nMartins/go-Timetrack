Usage:

Use case:
My personal use case was to be able to monitor a bunch of big function and decide
what changes have a big impact on runtime.
The scope is not to be able to improve the runtime down to Millisecond level.

To monitor execution of a function simply put
defer goTtrack.Timetrack(time.Now())
at the start of the function that you want to be tracked.

If you want to add timepoints in your function use Track instead.
For creating a timepoint just call TimePoint

This is my first real Project so feedback is very much welcome!
