// Package calendarops contains deterministic date-only calendar and scheduling rules.
package calendarops

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"
	"time"
)

// Date is a Gregorian civil date. It is not an instant and has no time zone.
type Date struct { Year int `json:"year"`; Month int `json:"month"`; Day int `json:"day"` }
func NewDate(y,m,d int) (Date,error) {
	if y < 1 || y > 9998 || m < 1 || m > 12 || d < 1 || d > DaysInMonth(y,m) { return Date{}, fmt.Errorf("invalid civil date %04d-%02d-%02d",y,m,d) }
	return Date{y,m,d},nil
}
func (d Date) String() string { return fmt.Sprintf("%04d-%02d-%02d",d.Year,d.Month,d.Day) }
func (d Date) time() time.Time { return time.Date(d.Year,time.Month(d.Month),d.Day,12,0,0,0,time.UTC) }
func FromTime(t time.Time) Date { y,m,d:=t.Date(); return Date{y,int(m),d} }
func ParseDate(s string) (Date,error) { t,e:=time.Parse("2006-01-02",s); if e!=nil { return Date{},e }; return NewDate(t.Year(),int(t.Month()),t.Day()) }
func (d Date) Weekday() time.Weekday { return d.time().Weekday() }
func (d Date) AddDays(n int) Date { return FromTime(d.time().AddDate(0,0,n)) }
func (d Date) Before(o Date) bool { return d.time().Before(o.time()) }
// DaysUntil returns signed civil-day distance from d to other.
func (d Date) DaysUntil(other Date) int { return int(other.time().Sub(d.time()).Hours()/24) }
func (d Date) Equal(o Date) bool { return d==o }
func IsLeapYear(y int) bool { return y%400==0 || (y%4==0 && y%100!=0) }
func DaysInMonth(y,m int) int { if m<1||m>12{return 0}; return time.Date(y,time.Month(m)+1,0,12,0,0,0,time.UTC).Day() }
func DayOfYear(d Date) int { return d.time().YearDay() }
func ISOWeek(d Date) (int,int) { return d.time().ISOWeek() }

type Cell struct { Date Date `json:"date"`; InMonth bool `json:"inMonth"`; ISOWeek int `json:"isoWeek"`; ISOYear int `json:"isoYear"` }
type Month struct { Year int `json:"year"`; Month int `json:"month"`; WeekStart string `json:"weekStart"`; Rows [][]Cell `json:"rows"` }
func MonthGrid(y,m int, weekStart string, fixed bool) (Month,error) {
	if y<1||y>9998||m<1||m>12{return Month{},fmt.Errorf("year must be 1..9998 and month 1..12")}
	if weekStart!="monday"&&weekStart!="sunday" {return Month{},fmt.Errorf("weekStart must be monday or sunday")}
	first:=Date{y,m,1}; wd:=int(first.Weekday()); start:=0; if weekStart=="monday"{start=1}; lead:=(wd-start+7)%7
	rows:=(lead+DaysInMonth(y,m)+6)/7; if fixed {rows=6}
	grid:=Month{Year:y,Month:m,WeekStart:weekStart,Rows:make([][]Cell,rows)}
	cursor:=first.AddDays(-lead)
	for r:=0;r<rows;r++ {grid.Rows[r]=make([]Cell,7); for c:=0;c<7;c++ {w,iy:=ISOWeek(cursor); grid.Rows[r][c]=Cell{cursor,cursor.Month==m,w,iy}; cursor=cursor.AddDays(1)}}
	return grid,nil
}
func YearMonths(y int, weekStart string) ([]Month,error) { if y<1||y>9998{return nil,fmt.Errorf("year must be 1..9998")}; out:=make([]Month,12); for m:=1;m<=12;m++ {g,e:=MonthGrid(y,m,weekStart,false);if e!=nil{return nil,e};out[m-1]=g};return out,nil }

// USFederalHolidays returns observed federal holidays. It models the current statutory
// set and labels cross-year observance by the civil date on which offices close.
type Holiday struct { Date Date `json:"date"`; Name string `json:"name"`; Kind string `json:"kind"` }
func nthWeekday(y,m int,w time.Weekday,n int) Date { d:=Date{y,m,1}; delta:=(int(w)-int(d.Weekday())+7)%7; return d.AddDays(delta+7*(n-1)) }
func lastWeekday(y,m int,w time.Weekday) Date { d:=Date{y,m,DaysInMonth(y,m)}; return d.AddDays(-((int(d.Weekday())-int(w)+7)%7)) }
func observed(d Date) Date { switch d.Weekday(){case time.Saturday:return d.AddDays(-1);case time.Sunday:return d.AddDays(1)};return d }
func USFederalHolidays(y int) []Holiday {
	if y < 2021 || y > 9998 { return nil }
	// New Year's Day is attributed by observed date, so the preceding year's observance appears here.
	items:=[]Holiday{{observed(Date{y,1,1}),"New Year's Day","federal"},{nthWeekday(y,1,time.Monday,3),"Birthday of Martin Luther King, Jr.","federal"},{nthWeekday(y,2,time.Monday,3),"Washington's Birthday","federal"},{lastWeekday(y,5,time.Monday),"Memorial Day","federal"},{Date{y,6,19},"Juneteenth National Independence Day","federal"},{Date{y,7,4},"Independence Day","federal"},{nthWeekday(y,9,time.Monday,1),"Labor Day","federal"},{nthWeekday(y,10,time.Monday,2),"Columbus Day","federal"},{Date{y,11,11},"Veterans Day","federal"},{nthWeekday(y,11,time.Thursday,4),"Thanksgiving Day","federal"},{Date{y,12,25},"Christmas Day","federal"}}
	for i:=range items {if items[i].Name!="New Year's Day" {items[i].Date=observed(items[i].Date)}}
	// Include next year's New Year observed on Dec 31 where applicable.
	next:=observed(Date{y+1,1,1}); if next.Year==y {items=append(items,Holiday{next,"New Year's Day (observed)","federal"})}
	sort.Slice(items,func(i,j int)bool{return items[i].Date.Before(items[j].Date)}); return items
}
func HolidaySet(years []int) map[Date]string { out:=map[Date]string{};for _,y:=range years{for _,h:=range USFederalHolidays(y){out[h.Date]=h.Name}};return out }
func HolidaySetBetween(start,end Date) map[Date]string { if end.Before(start){return map[Date]string{}};ys:=make([]int,0,end.Year-start.Year+1);for y:=start.Year;y<=end.Year;y++{ys=append(ys,y)};return HolidaySet(ys) }

type BusinessOptions struct { Weekends []time.Weekday; Excluded map[Date]string }
func DefaultBusinessOptions() BusinessOptions {return BusinessOptions{[]time.Weekday{time.Saturday,time.Sunday},map[Date]string{}}}
func (o BusinessOptions) IsBusinessDay(d Date) bool {for _,w:=range o.Weekends{if d.Weekday()==w{return false}};_,excluded:=o.Excluded[d];return !excluded}
// CountBusinessDays counts dates strictly after start through end, unless includeStart is true.
func CountBusinessDays(start,end Date,includeStart,includeEnd bool,o BusinessOptions) (int,error) {if end.Before(start){return 0,fmt.Errorf("end must be on or after start")};n:=0;d:=start;if !includeStart{d=d.AddDays(1)};for d.Before(end) || (includeEnd && d.Equal(end)) {if o.IsBusinessDay(d){n++};d=d.AddDays(1)};return n,nil}
func AddBusinessDays(anchor Date,n int,o BusinessOptions) (Date,error) {if n<0{return Date{},fmt.Errorf("business-day offset must be nonnegative")};d:=anchor;for left:=n;left>0;{d=d.AddDays(1);if o.IsBusinessDay(d){left--};if d.Year>9998{return Date{},fmt.Errorf("result exceeds supported year")}};return d,nil}
func NextBusinessDay(d Date,o BusinessOptions) Date { for {d=d.AddDays(1);if o.IsBusinessDay(d){return d}} }
func PreviousBusinessDay(d Date,o BusinessOptions) Date { for {d=d.AddDays(-1);if o.IsBusinessDay(d){return d}} }
func SubtractBusinessDays(anchor Date,n int,o BusinessOptions) (Date,error) {if n<0{return Date{},fmt.Errorf("business-day offset must be nonnegative")};d:=anchor;for left:=n;left>0;{d=d.AddDays(-1);if o.IsBusinessDay(d){left--};if d.Year<1{return Date{},fmt.Errorf("result precedes supported year")}};return d,nil}

// NRF445 anchors a fiscal year to the Sunday nearest February 1. Pattern repeats
// each quarter; an extra week, when present, extends period 12.
type FiscalPeriod struct { FiscalYear int `json:"fiscalYear"`; Period int `json:"period"`; Quarter int `json:"quarter"`; Start Date `json:"start"`; End Date `json:"end"`; Weeks int `json:"weeks"` }
func NRF445Start(y int) Date { feb:=Date{y,2,1}; sunday:=feb.AddDays(-int(feb.Weekday())); next:=sunday.AddDays(7); if feb.time().Sub(sunday.time()) <= next.time().Sub(feb.time()){return sunday}; return next }
func FiscalYearWeeks(y int) int { return int(NRF445Start(y+1).time().Sub(NRF445Start(y).time()).Hours()/168) }
func RetailPeriods(y int,pattern string) ([]FiscalPeriod,error) {var lens []int;switch pattern{case "445":lens=[]int{4,4,5};case "454":lens=[]int{4,5,4};case "544":lens=[]int{5,4,4};default:return nil,fmt.Errorf("pattern must be 445, 454, or 544")};weeks:=FiscalYearWeeks(y); if weeks!=52&&weeks!=53{return nil,fmt.Errorf("invalid anchored fiscal year length")};out:=make([]FiscalPeriod,12);cur:=NRF445Start(y);for p:=1;p<=12;p++{w:=lens[(p-1)%3];if p==12&&weeks==53{w++};end:=cur.AddDays(w*7-1);out[p-1]=FiscalPeriod{y,p,(p-1)/3+1,cur,end,w};cur=end.AddDays(1)};return out,nil }

type PayPeriod struct { Nominal Date `json:"nominalPayDate"`; Adjusted Date `json:"adjustedPayDate"` }
func monthlyStep(d Date,day int) Date { t:=time.Date(d.Year,time.Month(d.Month)+1,1,12,0,0,0,time.UTC); if day>DaysInMonth(t.Year(),int(t.Month())){day=DaysInMonth(t.Year(),int(t.Month()))}; return Date{t.Year(),int(t.Month()),day} }
func PayDates(start,end,anchor Date,cadence,adjust string,o BusinessOptions) ([]PayPeriod,error) {if end.Before(start){return nil,fmt.Errorf("end must be on or after start")};var out []PayPeriod; emit:=func(d Date){if !d.Before(start)&&!end.Before(d){a:=d;for !o.IsBusinessDay(a){if adjust=="previous"{a=a.AddDays(-1)}else{a=a.AddDays(1)}};out=append(out,PayPeriod{d,a})}}
	switch cadence {case "weekly","biweekly":step:=7;if cadence=="biweekly"{step=14};d:=anchor;for d.Before(start){d=d.AddDays(step)};for !end.Before(d){emit(d);d=d.AddDays(step)}
	case "monthly":d:=anchor;for d.Before(start){d=monthlyStep(d,anchor.Day)};for !end.Before(d){emit(d);d=monthlyStep(d,anchor.Day)}
	case "semimonthly":if anchor.Day>28{return nil,fmt.Errorf("semi-monthly anchor day must be 1..28")};for m:=time.Date(start.Year,time.Month(start.Month),1,12,0,0,0,time.UTC);!(Date{m.Year(),int(m.Month()),1}).Before(Date{end.Year,end.Month,1});m=m.AddDate(0,1,0){days:=[]int{anchor.Day};if anchor.Day!=15{days=append(days,15)};for _,day:=range days{emit(Date{m.Year(),int(m.Month()),day})}}
	default:return nil,fmt.Errorf("cadence must be weekly, biweekly, semimonthly, or monthly")};return out,nil }
func ShiftOn(d,anchor Date,cycle []bool) (bool,error) {if len(cycle)==0{return false,fmt.Errorf("shift cycle cannot be empty")};delta:=int(d.time().Sub(anchor.time()).Hours()/24);idx:=((delta%len(cycle))+len(cycle))%len(cycle);return cycle[idx],nil}
func NamedShift(name string) ([]bool,error) {switch name{case "4-on-3-off":return []bool{true,true,true,true,false,false,false},nil;case "5-on-2-off":return []bool{true,true,true,true,true,false,false},nil;case "2-2-3":return []bool{true,true,false,false,true,true,true,false,false,true,true,false,false,false},nil;default:return nil,fmt.Errorf("shift must be 4-on-3-off, 5-on-2-off, or 2-2-3")}}

func RenderMonthSVG(g Month,holidays map[Date]string) string {const width=900; const cellW=width/7;var b strings.Builder;fmt.Fprintf(&b,"<svg xmlns=\"http://www.w3.org/2000/svg\" role=\"img\" aria-label=\"%s %d calendar\" viewBox=\"0 0 900 670\"><rect width=\"100%%\" height=\"100%%\" fill=\"#fffdf8\"/><text x=\"24\" y=\"44\" font-family=\"sans-serif\" font-size=\"28\" font-weight=\"700\" fill=\"#14283f\">%s %d</text>",time.Month(g.Month).String(),g.Year,time.Month(g.Month).String(),g.Year); labels:=[]string{"Sun","Mon","Tue","Wed","Thu","Fri","Sat"};if g.WeekStart=="monday"{labels=[]string{"Mon","Tue","Wed","Thu","Fri","Sat","Sun"}};for i,s:=range labels{fmt.Fprintf(&b,"<text x=\"%d\" y=\"82\" font-family=\"sans-serif\" font-size=\"14\" fill=\"#536b80\">%s</text>",i*cellW+12,html.EscapeString(s))};rh:=82;for r,row:=range g.Rows{for c,cell:=range row{x:=c*cellW;y:=100+r*rh;fill:="#ffffff";if !cell.InMonth{fill="#f2f5f4"};if _,ok:=holidays[cell.Date];ok{fill="#d9f1ec"};fmt.Fprintf(&b,"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"6\" fill=\"%s\" stroke=\"#dce5e3\"/>",x,y,cellW-5,rh-5,fill);col:="#14283f";if !cell.InMonth{col="#8999a4"};fmt.Fprintf(&b,"<text x=\"%d\" y=\"%d\" font-family=\"sans-serif\" font-size=\"15\" fill=\"%s\">%d</text>",x+12,y+24,col,cell.Date.Day);if name,ok:=holidays[cell.Date];ok{if len(name)>18{name=name[:17]+"…"};fmt.Fprintf(&b,"<text x=\"%d\" y=\"%d\" font-family=\"sans-serif\" font-size=\"10\" fill=\"#176d67\">%s</text>",x+12,y+44,html.EscapeString(name))}}};b.WriteString("</svg>");return b.String() }
func WriteCSV(w io.Writer,rows []PayPeriod) error {c:=csv.NewWriter(w);if e:=c.Write([]string{"nominal_pay_date","adjusted_pay_date"});e!=nil{return e};for _,p:=range rows{if e:=c.Write([]string{p.Nominal.String(),p.Adjusted.String()});e!=nil{return e}};c.Flush();return c.Error()}
func EscapeICS(s string) string {s=strings.ReplaceAll(s,"\\","\\\\");s=strings.ReplaceAll(s,";","\\;");s=strings.ReplaceAll(s,",","\\,");s=strings.ReplaceAll(s,"\n","\\n");return s}
func ICSHolidayEvents(holidays []Holiday) string {var b strings.Builder;b.WriteString("BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//CalendarOps Suite//EN\r\nCALSCALE:GREGORIAN\r\n");for i,h:=range holidays{fmt.Fprintf(&b,"BEGIN:VEVENT\r\nUID:%04d-%02d-%02d-%d@calendarops\r\nDTSTAMP:%s\r\nDTSTART;VALUE=DATE:%04d%02d%02d\r\nDTEND;VALUE=DATE:%04d%02d%02d\r\nSUMMARY:%s\r\nEND:VEVENT\r\n",h.Date.Year,h.Date.Month,h.Date.Day,i,h.Date.time().Format("20060102T000000Z"),h.Date.Year,h.Date.Month,h.Date.Day,h.Date.AddDays(1).Year,h.Date.AddDays(1).Month,h.Date.AddDays(1).Day,EscapeICS(h.Name))};b.WriteString("END:VCALENDAR\r\n");return b.String()}
