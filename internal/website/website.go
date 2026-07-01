package website

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wbc-calendar/internal/calendar"
	"wbc-calendar/internal/config"
)

type Website struct {
	Tournaments  []*calendar.Tournament
	Tournaments1 []*calendar.Tournament
	Tournaments2 []*calendar.Tournament
	Calendars    map[string]*calendar.Calendar
	Others       []*calendar.OtherSchedule

	WBCSite string
	Title   string
	Updated time.Time
}

func CreateWebsite(s *calendar.Schedule, config *config.Config) {
	if err := copyStaticFiles(config.OutputDirectory); err != nil {
		log.Printf("Error copying website resources: %v", err)
	}

	tournaments := s.SortedTournamentList()
	var l int = len(tournaments) / 2

	site := Website{
		Tournaments:  tournaments,
		Tournaments1: tournaments[0:l],
		Tournaments2: tournaments[l:],
		Calendars:    s.Calendars,
		Others:       s.Others,

		WBCSite: fmt.Sprintf("http://boardgamers.org/wbc%02d/schedule.html", config.Year%100),
		Title:   fmt.Sprintf("WBC %d Event Schedule", config.Year),
		Updated: time.Now(),
	}

	if err := WriteSiteFiles(site, config); err != nil {
		log.Printf("Error rendering website template: %v", err)
	}
}

func copyStaticFiles(outputDir string) error {
	staticDir := filepath.Join("web", "static")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	return filepath.Walk(staticDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(staticDir, srcPath)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		dstPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		return copyFile(srcPath, dstPath)
	})
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return dst.Sync()
}

func WriteSiteFiles(site Website, config *config.Config) error {
	tmpl, err := template.ParseGlob("web/templates/*.gohtml")
	if err != nil {
		return err
	}

	if err := WriteTemplate(tmpl, site, config.OutputDirectory, "index", "index.html"); err != nil {
		return err
	}

	reportName := strings.Split(config.ExcelFilePath, "/")[len(strings.Split(config.ExcelFilePath, "/"))-1]
	reportName = strings.TrimSuffix(reportName, filepath.Ext(reportName))
	reportName = fmt.Sprintf("report-%s.html", reportName)
	if err := WriteTemplate(tmpl, site, config.OutputDirectory, "report", reportName); err != nil {
		return err
	}

	log.Printf("Website generated at %s", filepath.Join(config.OutputDirectory, "index.html"))

	return nil
}

func WriteTemplate(tmpl *template.Template, data interface{}, outputDir string, templateName string, outputName string) error {
	outPath := filepath.Join(outputDir, outputName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := tmpl.ExecuteTemplate(outFile, templateName, data); err != nil {
		return err
	}

	return outFile.Sync()
}
