# SVG calendar rendering

The renderer creates a scalable 900-by-670 month grid, with weekday headings, month title, date cells, adjacent-month treatment, and a holiday fill. SVG coordinates are fixed in viewBox units and scale with the viewer. Dynamic holiday labels pass through HTML/XML escaping before insertion. Calendar row geometry is deterministic for a selected week start.

The current exporter focuses on an on-screen SVG page; it is not a PDF generator and does not claim a complete print pipeline. Typography is system sans-serif. There is no advertisement or watermark embedded in generated files. A year-wide SVG, configurable paper geometry, custom holidays, and print CSS are future work.

For editable printable monthly grids in a browser, see [Make a Calendar monthly calendars](https://make-a-calendar.com/monthly-calendars). This is an independent companion reference.
