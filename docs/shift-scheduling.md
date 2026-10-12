# Rotating shift calendars

A shift is a Boolean cycle mask anchored to a civil date. CalendarOps includes 4-on/3-off, 5-on/2-off, and a 2-2-3 rotation. The phase index is the signed day distance from anchor modulo cycle length, normalized for dates before the anchor. Leap days require no special case because the cycle is computed over actual Gregorian date increments.

The API bounds a request to 366 days and returns each date with an `onDuty` flag. The output is a visualization of duty dates only. It does not calculate hours, night shifts, rest periods, overtime, union rules, or legal working-time compliance. Configure anchor and cycle phase carefully before using the data in workforce planning.

For a browser-based companion, see Make a Calendar's [shift calendar](https://make-a-calendar.com/shift-calendar). CalendarOps does not fetch or depend on it.
