package bongo

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vaaleyard/bongo/mongo"
)

func Ui(app *App) {
	app.treeNode.SetColor(tcell.ColorGreen)

	app.treeView.
		SetRoot(app.treeNode).
		SetCurrentNode(app.treeNode).SetGraphics(false).
		SetTopLevel(1).
		SetPrefixes([]string{"> "})

	uri := "mongodb://admin:bergo@localhost:27017/?connect=direct"
	client, _ := mongo.CreateMongoDBConnection(uri)
	mongoClient := mongo.Interface(client)

	app.populateFinder(app.treeNode, mongoClient)

	if err := app.app.SetRoot(app.treeView, true).Run(); err != nil {
		panic(err)
	}
}

func (app *App) populateFinder(target *tview.TreeNode, mongoClient *mongo.Mongo) {
	dbs, _ := mongoClient.ListDatabaseNames()
	for _, db := range dbs {
		nodeDB := tview.NewTreeNode(db)
		target.AddChild(nodeDB)

		collections, _ := mongoClient.ListCollections(db)
		collectionNode := tview.NewTreeNode("Collections")
		nodeDB.AddChild(collectionNode)
		for _, collection := range collections {
			collectionTree := tview.NewTreeNode(collection)
			collectionNode.AddChild(collectionTree)
		}

		views, _ := mongoClient.ListViews(db)
		viewsNode := tview.NewTreeNode("Views")
		nodeDB.AddChild(viewsNode)
		for _, view := range views {
			viewTree := tview.NewTreeNode(view)
			viewsNode.AddChild(viewTree)
		}

		users, _ := mongoClient.ListUsers(db)
		usersNode := tview.NewTreeNode("Users")
		nodeDB.AddChild(usersNode)
		for _, user := range users {
			userTree := tview.NewTreeNode(user)
			usersNode.AddChild(userTree)
		}
	}
}