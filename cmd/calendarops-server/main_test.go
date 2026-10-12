package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)
func TestCalendarAndErrorRoutes(t *testing.T){s:=&server{};h:=s.routes();tests:=[]struct{path string;status int;contains string}{{"/healthz",200,"ok"},{"/readyz",200,"ready"},{"/v1/calendar/month/2027/1?weekStart=monday",200,"2027"},{"/v1/calendar/month/2027/13",400,"error"},{"/v1/fiscal/445?year=2027",200,"periods"},{"/v1/exports/holidays.ics?year=2027",200,"BEGIN:VCALENDAR"}};for _,tc:=range tests{r:=httptest.NewRequest("GET",tc.path,nil);w:=httptest.NewRecorder();h.ServeHTTP(w,r);if w.Code!=tc.status{t.Errorf("%s status=%d want %d (%s)",tc.path,w.Code,tc.status,w.Body.String())};if !strings.Contains(w.Body.String(),tc.contains){t.Errorf("%s response did not contain %q",tc.path,tc.contains)}}}
func TestBusinessRangeFlags(t *testing.T){s:=&server{};r:=httptest.NewRequest("GET","/v1/business-days?start=2027-01-01&end=2027-01-04&includeStart=true&includeEnd=true",nil);w:=httptest.NewRecorder();s.routes().ServeHTTP(w,r);if w.Code!=200{t.Fatalf("status=%d %s",w.Code,w.Body.String())};var body map[string]any;if e:=json.Unmarshal(w.Body.Bytes(),&body);e!=nil{t.Fatal(e)};if body["count"]!=float64(1){t.Fatalf("count=%v",body["count"])};if body["includeStart"]!=true||body["includeEnd"]!=true{t.Fatal(body)}}
