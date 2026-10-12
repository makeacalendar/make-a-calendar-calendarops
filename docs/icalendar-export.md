# Automated iCalendar exports with CronJobs

The holiday exporter creates date-only events with `DTSTART;VALUE=DATE` and an exclusive `DTEND;VALUE=DATE` on the next civil date. The subset emitted is appropriate for all-day holidays; it does not model timed events, recurrence rules, attendees, alarms, or time-zone transitions. Text escaping covers backslashes, commas, semicolons, and newlines. Current event summaries are short enough not to need RFC 5545 line folding; any future user-defined long content must add folding.

A Kubernetes CronJob can run the standalone CLI without contacting the API. Schedule time zone controls when a Job starts; it does not convert holiday dates into instants. The default job chart uses `emptyDir`, which is removed with the pod. Set `output.ephemeral=false` to create/use a PVC and choose access mode, storage class, backup, retention, and a reader workload for the generated files.

The browser [iCalendar generator](https://make-a-calendar.com/ical-generator) is a useful way to inspect a one-off calendar. It is not used at runtime.
