package bot

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// add appends an item to a specified list. If the list
// does not exist, a new list is created. This function
// returns the text that the bot should send to the server.
func (b *Bot) add(guildID string, args []string) string {
	// Check that there are enough arguments given
	if len(args) < 2 {
		return `
		Add Command Error - Incorrect number of arguments
		
		Usage: ~add <list> <item> [items]
		`
	}

	var resultStr string

	filter := bson.M{"guildID": guildID}
	for _, item := range args[1:] {
		// Create the item to be added to a list
		itemBSON := bson.M{"$push": bson.M{"lists." + args[0]: item}}

		result, err := b.collection.UpdateOne(
			context.TODO(),
			filter,
			itemBSON,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fmt.Printf("Error adding item - %v\n", err)
			resultStr += fmt.Sprintf("Error - Could not add %s\n", item)
		}

		// Check if the item was added
		if result.ModifiedCount != 0 {
			resultStr += fmt.Sprintf("Added: %v to %v\n", item, args[0])
		}

		// Check if a new list had to be created
		if result.UpsertedCount != 0 {
			resultStr += fmt.Sprintf("Created List: %v\nAdded: %v to %v\n", args[0], item, args[0])
		}
	}

	return resultStr
}

// removeItem takes away one instance of an item from a specified
// list for a server. This function returns the text that the bot
// should send to the server.
func (b *Bot) removeItem(guildID string, args []string) string {
	// Check enough arguments for the command were given
	if len(args) < 2 {
		return `
		Remove Command Error - Incorrect number of arguments

		Usage: ~remove <list> <item> [items]
		`
	}

	var resultStr string

	filter := bson.M{"guildID": guildID}
	for _, item := range args[1:] {
		// Create the item to be added to a list
		itemBSON := bson.M{"$pull": bson.M{"lists." + args[0]: item}}

		result, err := b.collection.UpdateOne(
			context.TODO(),
			filter,
			itemBSON,
		)
		if err != nil {
			fmt.Printf("Error removing item - %v\n", err)
			resultStr += fmt.Sprintf("Error - Could not remove %s\n", item)
		}

		// Check if the item was added
		if result.ModifiedCount != 0 {
			resultStr += fmt.Sprintf("Removed: %v from %v\n", item, args[0])
		} else {
			resultStr += fmt.Sprintf("Error - %v in %v was not found\n", item, args[0])
		}
	}

	return resultStr
}

// removeList deletes an entire list from the discord server.
// This function returns the text that the bot should send
// to the server.
func (b *Bot) removeList(guildID string, args []string) string {
	// Check enough arguments for the command were given
	if len(args) < 1 {
		return `
		Remove List Command Error - Incorrect number of arguments

		Usage: ~removeList <list>
		`
	}

	filter := bson.M{"guildID": guildID}
	listBSON := bson.M{"$unset": bson.M{"lists." + args[0]: ""}}

	result, err := b.collection.UpdateOne(
		context.TODO(),
		filter,
		listBSON,
	)
	if err != nil {
		fmt.Printf("Error removing list - %v\n", err)
		return fmt.Sprintf("Error - Could not remove %s\n", args[0])
	}

	// Check if the item was added
	if result.ModifiedCount != 0 {
		return fmt.Sprintf("Removed: %s\n", args[0])
	} else {
		return fmt.Sprintf("Error - %s was not found\n", args[0])
	}
}

// // pick chooses at random one item from a specified list.
// // An optional argument can be provided which will remove
// // the chosen item from the list. This function returns
// // the text that the bot should send to the server.
// func (b *Bot) pick(guildID string, args []string) string {
// 	// Check there are enough arguments
// 	if len(args) < 1 {
// 		return "Pick Command Error - Incorrect number of arguments\n\nUsage: ~pick <list> [remove]"
// 	}

// 	// Check the specified list exists
// 	if _, exists := ds.lists[args[0]]; !exists {
// 		return fmt.Sprintf("Pick Command Error - Could not find list: %s", args[0])
// 	}

// 	// Get a random number between 0 and the length of the list.
// 	// Then get the chosen item from the list.
// 	chosenIndex := rand.Intn(len(ds.lists[args[0]]))
// 	chosen := ds.lists[args[0]][chosenIndex]

// 	// Check the remove argument has been provided
// 	if len(args) >= 2 {
// 		// Remove the chosen item from the list
// 		ds.lists[args[0]] = RemoveIndex(ds.lists[args[0]], chosenIndex)
// 		return fmt.Sprintf("Random item from list %s chosen and removed: %s", args[0], chosen)
// 	}

// 	return fmt.Sprintf("Random item from list %s chosen: %s", args[0], chosen)
// }

// // getList returns the list of items in the specified list.
// func (b *Bot) getList(guildID string, args []string) string {
// 	// Check enough arguments were given
// 	if len(args) < 1 {
// 		return `
// 		List Command Error - Incorrect number of arguments

// 		Usage: ~list <list>
// 		`
// 	}

// 	if _, exists := ds.lists[args[0]]; !exists {
// 		return fmt.Sprintf("List Command Error - Could not find list: %s", args[0])
// 	}

// 	result := fmt.Sprintf("List: %s\n", args[0])

// 	if len(ds.lists[args[0]]) == 0 {
// 		result += "\tNo Items"
// 		return result
// 	}

// 	for _, item := range ds.lists[args[0]] {
// 		result += fmt.Sprintf("\t%s\n", item)
// 	}

// 	return result
// }

// // getLists returns all of the list names stored by the server.
// func (b *Bot) getLists(guildID string) string {
// 	result := "Lists:\n"

// 	if len(ds.lists) == 0 {
// 		result += "\tNo Lists"
// 		return result
// 	}

// 	for list := range ds.lists {
// 		result += fmt.Sprintf("\t%s\n", list)
// 	}

// 	return result
// }
