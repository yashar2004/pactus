//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
)

//go:embed assets/ui/main_window.ui
var uiMainWindow []byte

type mainWindow struct {
	*gtk.ApplicationWindow

	widgetNode   *widgetNode
	widgetWallet *widgetWallet
}

func buildMainWindow(nodeModel *nodeModel, walletModel *walletModel) *mainWindow {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	fatalErrorCheck(err)

	appWindow := getApplicationWindowObj(builder, "id_main_window")
	boxNode := getBoxObj(builder, "id_box_node")
	boxDefaultWallet := getBoxObj(builder, "id_box_default_wallet")

	widgetNode, err := buildWidgetNode(nodeModel)
	fatalErrorCheck(err)

	widgetWallet, err := buildWidgetWallet(walletModel)
	fatalErrorCheck(err)

	boxNode.Add(widgetNode)
	boxDefaultWallet.Add(widgetWallet)

	mw := &mainWindow{
		ApplicationWindow: appWindow,
		widgetNode:        widgetNode,
		widgetWallet:      widgetWallet,
	}

	explorerItemMenu := getMenuItem(builder, "id_explorer_menu")
	explorerItemMenu.Connect("activate", mw.onMenuItemActivateExplorer)

	websiteItemMenu := getMenuItem(builder, "id_website_menu")
	websiteItemMenu.Connect("activate", mw.onMenuItemActivateWebsite)

	learnItemMenu := getMenuItem(builder, "id_learn_menu")
	learnItemMenu.Connect("activate", mw.onMenuItemActivateLearn)

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_about_gtk":            mw.onAboutGtk,
		"on_about":                mw.onAbout,
		"on_quit":                 mw.onQuit,
		"on_transaction_transfer": mw.OnTransactionTransfer,
		"on_transaction_bond":     mw.OnTransactionBond,
	}
	builder.ConnectSignals(signals)

	return mw
}

func (mw *mainWindow) onQuit() {
	mw.Close()
}

func (mw *mainWindow) onAboutGtk() {
	showAboutGTKDialog()
}

func (mw *mainWindow) onAbout() {
	showAboutDialog()
}

func (mw *mainWindow) OnTransactionTransfer() {
	broadcastTransactionSend(mw.widgetWallet.model.wallet)
}

func (mw *mainWindow) onMenuItemActivateWebsite(_ *gtk.MenuItem) {
	if err := util.OpenURLInBrowser("https://pactus.org/"); err != nil {
		fatalErrorCheck(err)
	}
}

func (mw *mainWindow) onMenuItemActivateExplorer(_ *gtk.MenuItem) {
	if err := util.OpenURLInBrowser("https://pactusscan.com/"); err != nil {
		fatalErrorCheck(err)
	}
}

func (mw *mainWindow) onMenuItemActivateLearn(_ *gtk.MenuItem) {
	if err := util.OpenURLInBrowser("https://pactus.org/learn/"); err != nil {
		fatalErrorCheck(err)
	}
}

func (mw *mainWindow) OnTransactionBond() {
	valAddrs := []crypto.Address{}
	consMgr := mw.widgetNode.model.node.ConsManager()
	for _, inst := range consMgr.Instances() {
		valAddrs = append(valAddrs, inst.SignerKey().Address())
	}
	broadcastTransactionBond(mw.widgetWallet.model.wallet, valAddrs)
}
