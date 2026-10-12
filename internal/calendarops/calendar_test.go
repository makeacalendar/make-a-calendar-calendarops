package calendarops

import (
	"strings"
	"testing"
	"time"
)

func mustDate(t *testing.T,s string) Date { t.Helper(); d,e:=ParseDate(s); if e!=nil {t.Fatal(e)};return d }
func TestLeapYearsAndMonthLengths(t *testing.T) {
	for _,tc:=range []struct{year,days int;leap bool}{{2024,29,true},{2100,28,false},{2000,29,true}} {if got:=IsLeapYear(tc.year);got!=tc.leap{t.Errorf("IsLeapYear(%d)=%v",tc.year,got)};if got:=DaysInMonth(tc.year,2);got!=tc.days{t.Errorf("DaysInMonth(%d,2)=%d",tc.year,got)}}
}
func TestISOYearRollover(t *testing.T) {for _,tc:=range []struct{date string;week,year int}{{"2026-12-31",53,2026},{"2027-01-01",53,2026},{"2027-01-04",1,2027}} {w,y:=ISOWeek(mustDate(t,tc.date));if w!=tc.week||y!=tc.year{t.Errorf("%s ISO=%d/%d want %d/%d",tc.date,w,y,tc.week,tc.year)}}}
func TestMonthGridAlignment(t *testing.T) {g,e:=MonthGrid(2027,1,"monday",false);if e!=nil{t.Fatal(e)};if g.Rows[0][0].Date.String()!="2026-12-28"{t.Fatalf("unexpected leading date %s",g.Rows[0][0].Date)};if g.Rows[len(g.Rows)-1][6].Date.String()!="2027-01-31"{t.Fatalf("unexpected trailing date %s",g.Rows[len(g.Rows)-1][6].Date)};s,_:=MonthGrid(2027,1,"sunday",false);if s.Rows[0][0].Date.String()!="2026-12-27"{t.Fatalf("unexpected Sunday start %s",s.Rows[0][0].Date)}}
func TestBusinessRangeInclusionAndHolidays(t *testing.T) {o:=DefaultBusinessOptions();o.Excluded=HolidaySet([]int{2027});start:=mustDate(t,"2027-01-01");end:=mustDate(t,"2027-01-04");got,e:=CountBusinessDays(start,end,false,true,o);if e!=nil{t.Fatal(e)};if got!=1{t.Fatalf("(start,end] count=%d, want 1",got)};got,e=CountBusinessDays(start,end,true,true,o);if e!=nil{t.Fatal(e)};if got!=1{t.Fatalf("inclusive count=%d, want 1 because Jan 1 is a holiday",got)}}
func TestObservedNewYearCrossYear(t *testing.T) {hs:=USFederalHolidays(2021);found:=false;for _,h:=range hs{if h.Date.String()=="2021-12-31"&&strings.Contains(h.Name,"New Year"){found=true}};if !found{t.Fatal("Saturday Jan 1 2022 should be observed on Dec 31 2021")}}
func TestFederalHolidayCoverageBeginsIn2021(t *testing.T) {if got:=USFederalHolidays(2020);len(got)!=0{t.Fatalf("historical schedule is outside coverage: %+v",got)};if got:=USFederalHolidays(2021);len(got)<11{t.Fatalf("expected 2021 statutory set, got %d",len(got))}}
func TestAddAndSubtractBusinessDays(t *testing.T) {o:=DefaultBusinessOptions();o.Excluded=HolidaySetBetween(mustDate(t,"2026-12-01"),mustDate(t,"2027-01-15"));start:=mustDate(t,"2027-01-01");next,e:=AddBusinessDays(start,1,o);if e!=nil{t.Fatal(e)};if next.String()!="2027-01-04"{t.Fatalf("next business date=%s",next)};prev,e:=SubtractBusinessDays(next,1,o);if e!=nil{t.Fatal(e)};if prev.String()!="2026-12-31"{t.Fatalf("previous business date=%s",prev)}}
func TestRetailPeriodCoverageAnd53WeekPlacement(t *testing.T) {for _,y:=range []int{2023,2024,2025,2026,2027,2028}{ps,e:=RetailPeriods(y,"445");if e!=nil{t.Fatal(e)};if len(ps)!=12{t.Fatalf("%d periods",len(ps))};if ps[0].Start!=NRF445Start(y){t.Fatal("incorrect anchor")};weeks:=0;for _,p:=range ps{weeks+=p.Weeks};if weeks!=FiscalYearWeeks(y){t.Fatalf("year %d: periods contain %d weeks",y,weeks)};if weeks==53&&ps[11].Weeks!=6{t.Fatalf("53rd week must extend period 12 in %d",y)}}}
func TestPayrollMonthlyClampsAnchorDay(t *testing.T) {o:=DefaultBusinessOptions();rows,e:=PayDates(mustDate(t,"2024-01-01"),mustDate(t,"2024-03-31"),mustDate(t,"2024-01-31"),"monthly","next",o);if e!=nil{t.Fatal(e)};if len(rows)!=3||rows[1].Nominal.String()!="2024-02-29"||rows[2].Nominal.String()!="2024-03-31"{t.Fatalf("unexpected dates: %+v",rows)}}
func TestShiftCyclePhaseAcrossLeapDay(t *testing.T) {cycle,_:=NamedShift("4-on-3-off");anchor:=mustDate(t,"2024-02-27");for i,want:=range []bool{true,true,true,true,false,false,false,true}{got,e:=ShiftOn(anchor.AddDays(i),anchor,cycle);if e!=nil{t.Fatal(e)};if got!=want{t.Errorf("offset %d=%v want %v",i,got,want)}}}
func TestSVGEscapesHolidayText(t *testing.T) {g,_:=MonthGrid(2027,1,"monday",false);out:=RenderMonthSVG(g,map[Date]string{{2027,1,1}:"A <B> & C"});if strings.Contains(out,"A <B>")||!strings.Contains(out,"A &lt;B&gt; &amp; C"){t.Fatal("holiday text was not escaped")}}
func TestICSDateOnlyAndExclusiveEnd(t *testing.T) {ics:=ICSHolidayEvents([]Holiday{{Date{2024,2,29},"Leap, day; special\\name","federal"}});for _,s:=range []string{"DTSTART;VALUE=DATE:20240229","DTEND;VALUE=DATE:20240301","SUMMARY:Leap\\, day\\; special\\\\name","BEGIN:VCALENDAR","END:VCALENDAR"}{if !strings.Contains(ics,s){t.Errorf("ICS missing %q",s)}};if !strings.Contains(ics,"DTSTAMP:"){t.Fatal("missing DTSTAMP")}}
func TestDateHasDateOnlySemantics(t *testing.T) {d:=mustDate(t,"2027-01-01");if d.Weekday()!=time.Friday{t.Fatal(d.Weekday())};if d.AddDays(1).String()!="2027-01-02"{t.Fatal("date arithmetic changed")}}

func TestTwoTwoThreeCycleIsFourteenDays(t *testing.T) {cycle,e:=NamedShift("2-2-3");if e!=nil{t.Fatal(e)};if len(cycle)!=14{t.Fatalf("cycle has %d days",len(cycle))};want:=[]bool{true,true,false,false,true,true,true,false,false,true,true,false,false,false};for i,v:=range want{if cycle[i]!=v{t.Errorf("phase %d=%v want %v",i,cycle[i],v)}}}
