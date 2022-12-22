package bongo

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vaaleyard/bongo/mongo"
)

func Ui(app *App) {
	treeNode := tview.NewTreeNode(".")
	treeNode.
		SetColor(tcell.ColorGreen)

	app.treeView.
		SetRoot(treeNode).
		SetCurrentNode(treeNode).SetGraphics(false).
		SetTopLevel(1).
		SetPrefixes([]string{"> "}).
		SetBorder(true).
		SetTitle("Finder").SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 0)

	app.inputBox.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).SetTitle("Input")

	app.preview.SetBorder(true).
		SetTitle("Preview").SetTitleAlign(tview.AlignCenter)

	layout := tview.NewFlex()
	layout.
		AddItem(app.inputBox, 0, 1, false).SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(app.treeView, 0, 1, true).
				AddItem(app.preview, 0, 5, false),
			0, 13, true)

	uri := "mongodb://admin:bergo@localhost:27017/?connect=direct"
	client, _ := mongo.CreateMongoDBConnection(uri)
	mongoClient := mongo.Interface(client)

	app.populateFinder(treeNode, mongoClient)
	app.treeView.SetSelectedFunc(selectNode)
	app.app.SetInputCapture(app.inputHandler)

	app.pages.AddPage("layout", layout, true, true)
	if err := app.app.SetRoot(app.pages, true).
		SetFocus(app.pages).Run(); err != nil {
		panic(err)
	}
}

func (app *App) populateFinder(target *tview.TreeNode, mongoClient *mongo.Mongo) {
	dbs, _ := mongoClient.ListDatabaseNames()
	for _, db := range dbs {
		nodeDB := tview.NewTreeNode(db)
		target.AddChild(nodeDB)

		collections, _ := mongoClient.ListCollections(db)
		collectionNode := tview.NewTreeNode("Collections").Collapse()
		nodeDB.AddChild(collectionNode)
		for _, collection := range collections {
			collectionTreeNode := tview.NewTreeNode(collection)
			collectionNode.AddChild(collectionTreeNode)
		}

		views, _ := mongoClient.ListViews(db)
		viewsNode := tview.NewTreeNode("Views").Collapse()
		nodeDB.AddChild(viewsNode)
		for _, view := range views {
			viewsTree := tview.NewTreeNode(view)
			viewsNode.AddChild(viewsTree)
		}

		users, _ := mongoClient.ListUsers(db)
		usersNode := tview.NewTreeNode("Users").Collapse()
		nodeDB.AddChild(usersNode)
		for _, user := range users {
			userTreeNode := tview.NewTreeNode(user)
			usersNode.AddChild(userTreeNode)
		}
	}
}

func selectNode(node *tview.TreeNode) {
	node.SetExpanded(!node.IsExpanded())
}

func (app *App) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 81 | 113 | 3: // Q or q or ctrl-c
		app.app.Stop()
		return nil
	case 73 | 105: //I or i
		app.app.SetFocus(app.inputBox)
	case 68 | 100: // D or d
		app.app.SetFocus(app.treeView)
	case 80 | 112: // P or p
		app.app.SetFocus(app.preview)
	}

	return event
}
