Usage:

Use case:
My personal use case was to be able to monitor a bunch of big functions and monitor
what changes have a noticable impact on runtime.

To monitor execution of a function simply put
`defer goTtrack.Timetrack(time.Now())`
at the start of the function that you want to be tracked.

If you want to add timepoints in your function use Track instead.
For creating a timepoint just call TimePoint

Use `GetCalcStatsPrint()` to get see your execution times.

This is my first real Project so feedback is very much welcome!
