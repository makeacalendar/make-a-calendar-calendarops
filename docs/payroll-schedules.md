# Payroll date schedules

The package generates nominal dates for weekly and biweekly cadence from an anchor date, monthly dates from the anchor day (clamped to the final day of a shorter month), and semi-monthly dates on the configured anchor day plus the 15th. Semi-monthly anchor days must be no later than 28. The requested range is inclusive of generated dates.

An adjustment policy reports a second date: `next` advances through non-business dates; `previous` moves backward. The nominal and adjusted values remain separately visible so a calendar does not hide the originally configured schedule. Default exclusions are weekends and current federal observed holidays. Custom calendars are supported by the library options but are not yet request parameters on the HTTP route.

CalendarOps does not calculate salary, tax, benefits, or a worker's legal entitlement. Payroll dates must be reviewed against applicable agreements and employer rules. The [pay-period calendar](https://make-a-calendar.com/pay-period-calendar) offers a complementary browser visualization.
