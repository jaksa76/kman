package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

var (
	l            *widgets.List
	instructions *widgets.Paragraph
	selection    []string
	filter       string
	root         string
	files        []string
)

func main() {
	root = "."
	files = listFiles(root)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	instructions = widgets.NewParagraph()
	instructions.Title = "Instructions"
	instructions.Text = "Start typing to filter files.     CTRL + r: restart     CTRL + a: apply     CTRL + d: delete     ESC: exit"

	l = widgets.NewList()
	l.Title = "Filter:"
	l.Rows = files
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	resize()
	ui.Render(l)
	ui.Render(instructions)

	updateSelection()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if isCharacter(e) {
			filter = filter + e.ID[0:1]
			updateSelection()
		}
		switch e.ID {
		case "<C-c>", "<Escape>":
			return
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<C-a>":
			applyAll(selection)
		case "<C-d>":
			deleteAll(selection)
		case "<C-r>":
			restartAll(selection)
		case "<Home>":
			l.ScrollTop()
		case "<End>":
			l.ScrollBottom()
		case "<Backspace>":
			if len(filter) > 0 {
				filter = filter[:len(filter)-1]
				updateSelection()
			}
		case "<Resize>":
			resize()
			ui.Render(instructions)
		}

		ui.Render(l)
	}
}

func isCharacter(e ui.Event) bool {
	return len(e.ID) == 1
	// return len(e.ID) == 1 && ((e.ID[0] >= 'a' && e.ID[0] <= 'z') ||
	// 	(e.ID[0] >= 'A' && e.ID[0] <= 'Z') ||
	// 	strings.IndexByte("*-._?", e.ID[0]) >= 0)
}

func resize() {
	w, _ := terminal.Width()
	h, _ := terminal.Height()
	l.SetRect(0, 0, int(w), int(h)-3)
	instructions.SetRect(0, int(h)-3, int(w), int(h))
}

func updateSelection() {
	filteredFiles, err := filterFiles(files, filter)
	if err == nil {
		selection = filteredFiles
		l.Title = "Filter: " + filter
		l.Rows = selection
	}
}

func filterFiles(files []string, filter string) (filtered []string, err error) {
	selection := make([]string, 0, len(files))

	re, err := regexp.Compile(toRegExp(filter))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if re.FindStringIndex(file) != nil {
			selection = append(selection, file)
		}
	}

	return selection, nil
}

func listFiles(root string) []string {
	var filenames []string

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return filenames
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func toRegExp(wildcard string) string {
	return strings.NewReplacer(
		".", "\\.",
		"?", ".",
		"*", ".*").Replace(wildcard)
}

func applyAll(selection []string) {
	ui.Close()
	for _, s := range selection {
		execWithOutput("kubectl", "apply", "-f", s)
	}
	os.Exit(0)
}

func deleteAll(selection []string) {
	ui.Close()
	for _, s := range selection {
		execWithOutput("kubectl", "delete", "-f", s)
	}
	os.Exit(0)

}

func restartAll(selection []string) {
	ui.Close()
	for _, s := range selection {
		execWithOutput("kubectl", "delete", "-f", s)
	}
	for _, s := range selection {
		execWithOutput("kubectl", "apply", "-f", s)
	}
	os.Exit(0)
}

func execWithOutput(command ...string) {
	fmt.Println(command)
	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
