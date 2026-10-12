# Date engine and civil calendar boundaries

A civil date identifies a day on the Gregorian calendar. `2027-01-01` is not midnight in a chosen time zone. CalendarOps stores an integer year, month, and day; standard-library UTC noon is used only as a safe arithmetic convention when asking Go for weekday or ISO-week facts. Serialized values remain date-only.

Leap-year rule: divisible by 400, or divisible by 4 and not 100. Therefore 2000 has February 29; 2100 does not. Month-grid offset is the first weekday minus the selected first weekday modulo seven. Rows use seven cells and include adjacent-month dates as needed; callers can request a fixed six-row layout in the library.

ISO week number is paired with ISO week-year. For example, January 1, 2027 belongs to ISO week 53 of ISO year 2026. The date engine uses the proleptic Gregorian rules and supports civil years 1–9998 for ordinary calendar operations. It does not model jurisdictional calendar reform history or locale-specific week-number systems.

The [monthly calendar pages](https://make-a-calendar.com/monthly-calendars) on Make a Calendar offer a complementary browser editor and printable layouts. The implementation here is independent and has no network dependency.
