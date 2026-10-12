# Building a 4-4-5 fiscal calendar engine

CalendarOps anchors fiscal year `Y` on the Sunday nearest February 1 in Gregorian year `Y`. The same policy is used to find the next fiscal year's start. An anchored year therefore contains 52 weeks or occasionally 53. The supported patterns repeat every quarter: 4-4-5, 4-5-4, and 5-4-4. These are thirteen-week retail quarters, not twelve ordinary Gregorian months.

When the anchor distance spans 53 weeks, CalendarOps extends period 12 by one week. Every period includes start and end civil dates, quarter number, and week count. The sum of period weeks equals the difference between consecutive anchors. `year` denotes the Gregorian year containing the February anchor; it is not a globally standardized fiscal-year naming convention.

This is an explicitly documented Sunday-nearest-February-1 policy, commonly described as NRF-style. Organizations may use a different year label or extra-week allocation; verify their convention before reporting financial results. The [4-4-5 planning tool](https://make-a-calendar.com/445-calendar) is a complementary browser reference.
