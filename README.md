# wbc-calendar
New Calendar tool for WBC

This tool is designed to read WBC Schedule data in an Excel file and produce the following:
* WebCal calendars for use with Calendar programs such as Google Calendars or Apple iCalendar.
* CSV and JSON formatted schedule data for the app developers.
* An HTML index easy use.
* A report of discrepancies in the schedule data to improve quality of the schedules.

# Usage
```
Usage: go run wbc-calendar [OPTIONS]

Excel spreadsheet reader for WBC Calendar project.

Options:
  -c    Clean output directory before writing (short form)
  -clean
        Clean output directory before writing
  -f string
        Path to Excel file (short form)
  -file string
        Path to Excel file (required)
  -l    List available sheets and exit (short form)
  -list
        List available sheets and exit
  -o string
        Output directory for calendar files (short form) (default "build")
  -output string
        Output directory for calendar files (default "build")
  -s string
        Sheet name to read (short form) (default "App Version")
  -sheet string
        Sheet name to read (default "App Version")
  -v    Enable verbose output (short form)
  -verbose
        Enable verbose output
  -y int
        Calendar year (short form)
  -year int
        Calendar year
  -z string
        Time zone location (short form) (default "America/New_York")
  -zone string
        Time zone location (default: America/New_York) (default "America/New_York")

Examples:
  go run wbc-calendar -y 2026 -f data.xlsx
  go run wbc-calendar -y 2026 -f data.xlsx -l
  go run wbc-calendar -y 2026 -f data.xlsx -s "Sheet1"
```
