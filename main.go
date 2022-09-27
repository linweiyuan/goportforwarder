package main

import (
	"github.com/linweiyuan/goportforwarder/k8s"
	"github.com/linweiyuan/goportforwarder/ui"
)

func main() {
	ui := ui.New()
	pods := k8s.GetPods()
	podTable := ui.NewPodTable(pods)
	ui.APP.SetRoot(podTable, true).Run()
}
