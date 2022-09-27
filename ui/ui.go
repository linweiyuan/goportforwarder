package ui

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/linweiyuan/goportforwarder/k8s"
	"github.com/rivo/tview"
)

type UI struct {
	APP *tview.Application
}

func New() *UI {
	return &UI{
		APP: tview.NewApplication(),
	}
}

type PodTable struct {
	*tview.Table
}

func (ui *UI) NewPodTable(pods []k8s.Pod) *PodTable {
	podTable := &PodTable{
		Table: tview.NewTable().SetSelectable(true, false).Select(1, 0),
	}

	podTable.SetTitle("Pods").SetBorder(true).SetTitleAlign(tview.AlignLeft)
	tableHeaders := []string{
		"Name",
		"Port",
	}

	for index, header := range tableHeaders {
		podTable.SetCell(0, index, tview.NewTableCell(header))
	}

	for index, pod := range pods {
		podTable.SetCell(index+1, 0, tview.NewTableCell(pod.Name))
		podTable.SetCell(index+1, 1, tview.NewTableCell(pod.Port))
	}

	podTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := podTable.GetSelection()
			if row == 0 {
				return event
			}
			pod := pods[row-1]
			podName := pod.Name
			podPort := pod.Port
			modal := tview.NewModal().SetText(fmt.Sprintf("Do you want to expose %s with port %s?", podName, podPort)).
				AddButtons([]string{"Yes", "No"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonLabel {
				case "Yes":
					exec.Command("sh", "-c", fmt.Sprintf("kubectl port-forward --address 0.0.0.0 %s %s:%s", podName, podPort, podPort)).Start()
					ui.APP.SetRoot(podTable, true).SetFocus(podTable)
				case "No":
					exec.Command("sh", "-c", fmt.Sprintf("netstat -tnpl | grep %s | awk '{print $7}' | awk -F/ '{print $1}' | xargs kill", podPort)).Start()
					ui.APP.SetRoot(podTable, true).SetFocus(podTable)
				}
			})
			ui.APP.SetRoot(modal, false).SetFocus(modal)
		}
		return event
	})
	return podTable
}
