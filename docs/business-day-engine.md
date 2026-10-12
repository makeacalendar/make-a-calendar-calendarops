# Business-day calculations at scale

The reusable engine accepts a weekend weekday set and a map of custom excluded dates. HTTP defaults are Saturday and Sunday plus modern US federal observed holidays. A custom company closure is conceptually separate from a federal holiday; employers should add it explicitly rather than assuming every organization closes for every federal holiday.

The API's default interval is `(start, end]`: the start date is excluded and end date included. Query flags `includeStart` and `includeEnd` change either endpoint independently. Reversed ranges return a client error. Adding zero business days preserves the anchor; a positive offset advances until that many business dates have been encountered. Subtraction is the reverse direction. These policies make month-end and inclusive payroll deadlines easier to specify without relying on an ambiguous “days between” label.

Weekend observations shift fixed-date federal holidays: Saturday to Friday and Sunday to Monday. When the following January 1 is observed on December 31, it is included with that civil date. This release exposes the current statutory set and does not reproduce earlier historical applicability dates before 2021; HTTP federal-holiday ranges must begin in 2021 or later. Do not use it as a historical legal calendar.

The browser-based [business-days calculator](https://make-a-calendar.com/business-days-calculator) is useful for an independent one-off date check. It is not an API dependency.
