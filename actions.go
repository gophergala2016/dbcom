package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gophergala2016/dbcom/db"
	"github.com/grafov/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func refresh(g *gocui.Gui, v *gocui.View) error {
	return g.Flush()
}

// cursorDown enlarge sql view when a new line added
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	_, sizeY := v.Size()
	if y >= sizeY-1 {
		if splitY == minSplit {
			return nil
		}
		splitY--
		l, _ := g.View("log")
		l.Resize(0, -1)
		v.Resize(0, 1)
		v.SetOrigin(0, 0)
		v.SetCursor(x, y+1)
	}
	return nil
}

func execQuery(g *gocui.Gui, v *gocui.View) error {
	dbc := db.Use("test")
	query := strings.TrimSpace(v.Buffer())
	v.Reset()
	fmt.Fprint(v, query)
	widtgh, height := v.Size()
	if height > v.LinesCount() {
		v.SetSize(widtgh, v.LinesCount())
		_, maxY := g.Size()
		splitY = maxY - v.LinesCount() - 3
	}
	v.SetCursor(0, 0)
	if len(query) == 0 {
		return nil
	}

	log, _ := g.View("log")
	rows, err := dbc.Queryx(query)
	if err != nil {
		fmt.Fprintln(log, err)
		return nil
	}

	fmt.Fprintln(log, "\n", query)
	logWidth, _ := log.Size()
	log.Write(bytes.Repeat([]byte("─"), logWidth))
	for rows.Next() {
		log.Write([]byte("\n"))
		cols, err := rows.SliceScan()
		if err != nil {
			fmt.Fprintln(log, err)
			err = nil
			break
		}
		for _, col := range cols {
			log.Write(col.([]byte))
			log.Write([]byte("\t"))
		}
	}
	log.Write([]byte("\n"))
	log.Write(bytes.Repeat([]byte("─"), logWidth))

	return err
}
