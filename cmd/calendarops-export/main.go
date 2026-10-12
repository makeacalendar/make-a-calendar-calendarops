package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	co "github.com/makeacalendar/make-a-calendar-calendarops/internal/calendarops"
)

type dayRow struct { Date string `json:"date"`; Weekday string `json:"weekday"`; InMonth bool `json:"inMonth"`; Holiday string `json:"holiday,omitempty"` }
func main() {
	kind:=flag.String("type","month","report: month, year, holidays, fiscal")
	format:=flag.String("format","json","output: json, csv, svg, ics")
	year:=flag.Int("year",2027,"Gregorian or fiscal year; 0 uses current UTC year")
	month:=flag.Int("month",1,"month number 1..12; 0 uses current UTC month")
	weekStart:=flag.String("week-start","monday","monday or sunday")
	pattern:=flag.String("pattern","445","retail fiscal pattern: 445, 454, 544")
	output:=flag.String("output","-","output file; '-' writes to stdout")
	startFlag:=flag.String("start","","business/payroll/shift start date YYYY-MM-DD")
	endFlag:=flag.String("end","","business/payroll/shift end date YYYY-MM-DD")
	anchorFlag:=flag.String("anchor","","payroll/shift anchor date YYYY-MM-DD")
	cadence:=flag.String("cadence","biweekly","payroll cadence")
	adjust:=flag.String("adjust","next","weekend/holiday adjustment: next or previous")
	includeStart:=flag.Bool("include-start",false,"include start endpoint in business-day calculation")
	includeEnd:=flag.Bool("include-end",true,"include end endpoint in business-day calculation")
	cycle:=flag.String("cycle","4-on-3-off","shift cycle")
	flag.Parse()
	if *year==0 { *year=time.Now().UTC().Year() }; if *month==0 { *month=int(time.Now().UTC().Month()) }; *output=strings.ReplaceAll(*output,"{year}",strconv.Itoa(*year)); *output=strings.ReplaceAll(*output,"{month}",fmt.Sprintf("%02d",*month))
	if *year<1||*year>9998 { fatal("year must be 1..9998") }
	var content []byte; var err error
	switch *kind {
	case "month": content,err=monthReport(*year,*month,*weekStart,*format)
	case "year": content,err=yearReport(*year,*weekStart,*format)
	case "holidays": content,err=holidayReport(*year,*format)
	case "fiscal": content,err=fiscalReport(*year,*pattern,*format)
	case "business-days": content,err=businessReport(*startFlag,*endFlag,*includeStart,*includeEnd,*format)
	case "pay-periods": content,err=payrollReport(*startFlag,*endFlag,*anchorFlag,*cadence,*adjust,*format)
	case "shift-schedule": content,err=shiftReport(*startFlag,*endFlag,*anchorFlag,*cycle,*format)
	default: fatal("type must be month, year, holidays, fiscal, business-days, pay-periods, or shift-schedule")
	}
	check(err)
	if *output=="-" { if _,err=os.Stdout.Write(content);err!=nil{fatal(err.Error())};return }
	if err=os.WriteFile(*output,content,0644);err!=nil{fatal(err.Error())}
}
func monthReport(y,m int,week,format string)([]byte,error){g,e:=co.MonthGrid(y,m,week,false);if e!=nil{return nil,e};switch format{case "json":return jsonBytes(g);case "svg":return []byte(co.RenderMonthSVG(g,holidayMap(y))),nil;case "csv":return rowsCSV(rowsForMonth(g,holidayMap(y)));default:return nil,fmt.Errorf("month format must be json, csv, or svg")}}
func yearReport(y int,week,format string)([]byte,error){months,e:=co.YearMonths(y,week);if e!=nil{return nil,e};hols:=holidayMap(y);var rows []dayRow;for _,m:=range months{rows=append(rows,rowsForMonth(m,hols)...)};switch format{case "json":return jsonBytes(map[string]any{"year":y,"months":months});case "csv":return rowsCSV(rows);default:return nil,fmt.Errorf("year format must be json or csv")}}
func holidayReport(y int,format string)([]byte,error){if y<2021{return nil,fmt.Errorf("holiday coverage begins in 2021")};hs:=co.USFederalHolidays(y);switch format{case "json":return jsonBytes(hs);case "ics":return []byte(co.ICSHolidayEvents(hs)),nil;case "csv":var rows []dayRow;for _,h:=range hs{rows=append(rows,dayRow{h.Date.String(),h.Date.Weekday().String(),true,h.Name})};return rowsCSV(rows);default:return nil,fmt.Errorf("holiday format must be json, csv, or ics")}}
func fiscalReport(y int,pattern,format string)([]byte,error){ps,e:=co.RetailPeriods(y,pattern);if e!=nil{return nil,e};switch format{case "json":return jsonBytes(map[string]any{"year":y,"pattern":pattern,"anchor":"Sunday nearest February 1","periods":ps});case "csv":var b bytes.Buffer;c:=csv.NewWriter(&b);_ = c.Write([]string{"fiscal_year","period","quarter","start","end","weeks"});for _,p:=range ps{_ = c.Write([]string{strconv.Itoa(p.FiscalYear),strconv.Itoa(p.Period),strconv.Itoa(p.Quarter),p.Start.String(),p.End.String(),strconv.Itoa(p.Weeks)})};c.Flush();return b.Bytes(),c.Error();default:return nil,fmt.Errorf("fiscal format must be json or csv")}}

func businessReport(startText,endText string,includeStart,includeEnd bool,format string)([]byte,error){start,e:=co.ParseDate(startText);if e!=nil{return nil,fmt.Errorf("start must be YYYY-MM-DD")};end,e:=co.ParseDate(endText);if e!=nil{return nil,fmt.Errorf("end must be YYYY-MM-DD")};if end.Before(start)||start.DaysUntil(end)>3660{return nil,fmt.Errorf("range must be ordered and at most ten years")};if start.Year<2021||end.Year>9998{return nil,fmt.Errorf("federal holiday range must be 2021..9998")};o:=co.DefaultBusinessOptions();o.Excluded=co.HolidaySetBetween(start,end);var rows []dayRow;for d:=start;!end.Before(d);d=d.AddDays(1){if d==start&&!includeStart||d==end&&!includeEnd{continue};if o.IsBusinessDay(d){rows=append(rows,dayRow{d.String(),d.Weekday().String(),true,""})}};if format=="csv"{return rowsCSV(rows)};if format=="json"{return jsonBytes(map[string]any{"start":start,"end":end,"includeStart":includeStart,"includeEnd":includeEnd,"businessDays":rows})};return nil,fmt.Errorf("business-days format must be json or csv")}
func payrollReport(startText,endText,anchorText,cadence,adjust,format string)([]byte,error){start,e:=co.ParseDate(startText);if e!=nil{return nil,fmt.Errorf("start must be YYYY-MM-DD")};end,e:=co.ParseDate(endText);if e!=nil{return nil,fmt.Errorf("end must be YYYY-MM-DD")};anchor,e:=co.ParseDate(anchorText);if e!=nil{return nil,fmt.Errorf("anchor must be YYYY-MM-DD")};if start.Year<2021||end.Year>9998||anchor.Year<2021{return nil,fmt.Errorf("federal holiday coverage is 2021..9998")};o:=co.DefaultBusinessOptions();lo:=start;if anchor.Before(lo){lo=anchor};hi:=end;if hi.Before(anchor){hi=anchor};o.Excluded=co.HolidaySetBetween(lo,hi);rows,e:=co.PayDates(start,end,anchor,cadence,adjust,o);if e!=nil{return nil,e};if format=="csv"{var b bytes.Buffer;e=co.WriteCSV(&b,rows);return b.Bytes(),e};if format=="json"{return jsonBytes(rows)};return nil,fmt.Errorf("pay-periods format must be json or csv")}
type shiftRow struct{Date string `json:"date"`;OnDuty bool `json:"onDuty"`}
func shiftReport(startText,endText,anchorText,cycleName,format string)([]byte,error){start,e:=co.ParseDate(startText);if e!=nil{return nil,fmt.Errorf("start must be YYYY-MM-DD")};end,e:=co.ParseDate(endText);if e!=nil{return nil,fmt.Errorf("end must be YYYY-MM-DD")};anchor,e:=co.ParseDate(anchorText);if e!=nil{return nil,fmt.Errorf("anchor must be YYYY-MM-DD")};if end.Before(start)||start.DaysUntil(end)>366{return nil,fmt.Errorf("range must be ordered and at most 366 days")};cycle,e:=co.NamedShift(cycleName);if e!=nil{return nil,e};var rows []shiftRow;for d:=start;!end.Before(d);d=d.AddDays(1){on,_:=co.ShiftOn(d,anchor,cycle);rows=append(rows,shiftRow{d.String(),on})};if format=="json"{return jsonBytes(map[string]any{"cycle":cycleName,"anchor":anchor,"days":rows})};if format=="csv"{var b bytes.Buffer;c:=csv.NewWriter(&b);_ = c.Write([]string{"date","on_duty"});for _,r:=range rows{_ = c.Write([]string{r.Date,strconv.FormatBool(r.OnDuty)})};c.Flush();return b.Bytes(),c.Error()};return nil,fmt.Errorf("shift-schedule format must be json or csv")}

func rowsForMonth(g co.Month,hols map[co.Date]string)[]dayRow{var rows []dayRow;for _,week:=range g.Rows{for _,cell:=range week{h:=hols[cell.Date];rows=append(rows,dayRow{cell.Date.String(),cell.Date.Weekday().String(),cell.InMonth,h})}};return rows}
func holidayMap(y int)map[co.Date]string{m:=map[co.Date]string{};for _,h:=range co.USFederalHolidays(y){m[h.Date]=h.Name};return m}
func rowsCSV(rows []dayRow)([]byte,error){var b bytes.Buffer;c:=csv.NewWriter(&b);if e:=c.Write([]string{"date","weekday","in_month","holiday"});e!=nil{return nil,e};for _,r:=range rows{if e:=c.Write([]string{r.Date,r.Weekday,strconv.FormatBool(r.InMonth),r.Holiday});e!=nil{return nil,e}};c.Flush();return b.Bytes(),c.Error()}
func jsonBytes(v any)([]byte,error){b,e:=json.MarshalIndent(v,"","  ");if e!=nil{return nil,e};return append(b,'\n'),nil}
func check(e error){if e!=nil{fatal(e.Error())}}
func fatal(s string){fmt.Fprintln(os.Stderr,"calendarops-export:",s);os.Exit(2)}
