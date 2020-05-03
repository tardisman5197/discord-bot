package bot

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// add appends an item to a specified list. If the list
// does not exist, a new list is created. This function
// returns the text that the bot should send to the server.
func (b *Bot) add(guildID string, args []string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "add",
		"args":    args,
	})

	// Check that there are enough arguments given
	if len(args) < 2 {
		logger.Debug("Invalid number of arguments")

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

		ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
		result, err := b.collection.UpdateOne(
			ctx,
			filter,
			itemBSON,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			logger.WithError(err).Error("Error adding item")
			resultStr += fmt.Sprintf("Error - Could not add %s\n", item)
		}

		// Check if the item was added
		if result.ModifiedCount != 0 {
			logger.Debug("Added item to list")
			resultStr += fmt.Sprintf("Added: %v to %v\n", item, args[0])
		}

		// Check if a new list had to be created
		if result.UpsertedCount != 0 {
			logger.Debug("Created new list and added item")
			resultStr += fmt.Sprintf("Created List: %v\nAdded: %v to %v\n", args[0], item, args[0])
		}
	}

	return resultStr
}

// removeItem takes away one instance of an item from a specified
// list for a server. This function returns the text that the bot
// should send to the server.
func (b *Bot) removeItem(guildID string, args []string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "remove",
		"args":    args,
	})

	// Check enough arguments for the command were given
	if len(args) < 2 {
		logger.Debug("Invalid number of arguments")

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

		ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
		result, err := b.collection.UpdateOne(
			ctx,
			filter,
			itemBSON,
		)
		if err != nil {
			logger.WithError(err).Error("Error removing item")
			resultStr += fmt.Sprintf("Error - Could not remove %s\n", item)
		}

		// Check if the item was added
		if result.ModifiedCount != 0 {
			logger.Debug("Removed item from list")
			resultStr += fmt.Sprintf("Removed: %v from %v\n", item, args[0])
		} else {
			logger.Debug("Could not remove item - no item found")
			resultStr += fmt.Sprintf("Error - %v in %v was not found\n", item, args[0])
		}
	}

	return resultStr
}

// removeList deletes an entire list from the discord server.
// This function returns the text that the bot should send
// to the server.
func (b *Bot) removeList(guildID string, args []string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "removeList",
		"args":    args,
	})

	// Check enough arguments for the command were given
	if len(args) < 1 {
		logger.Debug("Invalid number of arguments")

		return `
		Remove List Command Error - Incorrect number of arguments

		Usage: ~removeList <list>
		`
	}

	filter := bson.M{"guildID": guildID}
	listBSON := bson.M{"$unset": bson.M{"lists." + args[0]: ""}}

	ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
	result, err := b.collection.UpdateOne(
		ctx,
		filter,
		listBSON,
	)
	if err != nil {
		logger.WithError(err).Error("Error removing list")
		return fmt.Sprintf("Error - Could not remove %s\n", args[0])
	}

	// Check if the item was added
	if result.ModifiedCount != 0 {
		logger.Debug("Removed list")
		return fmt.Sprintf("Removed: %s\n", args[0])
	}

	logger.Debug("Could not remove list - no list found")
	return fmt.Sprintf("Error - %s was not found\n", args[0])
}

// pick chooses at random one item from a specified list.
// An optional argument can be provided which will remove
// the chosen item from the list. This function returns
// the text that the bot should send to the server.
func (b *Bot) pick(guildID string, args []string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "pick",
		"args":    args,
	})

	// Check there are enough arguments
	if len(args) < 1 {
		logger.Debug("Invalid number of arguments")
		return "Pick Command Error - Incorrect number of arguments\n\nUsage: ~pick <list> [remove]"
	}

	// Get the list
	var result bson.M
	ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
	err := b.collection.FindOne(
		ctx,
		bson.M{"guildID": guildID, "lists." + args[0]: bson.M{"$exists": true}},
	).Decode(&result)

	if err != nil {
		logger.WithError(err).Error("Error finding list")
		return fmt.Sprintf("Error finding list - %v", err)
	}

	items := result["lists"].(bson.M)[args[0]].(bson.A)
	if len(items) == 0 {
		logger.Debug("Could not pick item - list empty")
		return fmt.Sprintf("No items in %v", args[0])
	}

	chosen := items[rand.Intn(len(items))]

	// Remove the item from the list
	var removedStr string
	if len(args) > 1 {
		removedStr += fmt.Sprintf("\n%s", b.removeItem(guildID, []string{args[0], fmt.Sprintf("%v", chosen)}))
	}

	logger.Debug("Picked an item")
	return fmt.Sprintf("Item chosen from %s: %v%s", args[0], chosen, removedStr)
}

// getList returns the list of items in the specified list.
func (b *Bot) getList(guildID string, args []string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "list",
		"args":    args,
	})

	// Check enough arguments were given
	if len(args) < 1 {
		logger.Debug("Invalid number of arguments")

		return `
		List Command Error - Incorrect number of arguments

		Usage: ~list <list>
		`
	}

	// Get the list
	var result bson.M
	ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
	err := b.collection.FindOne(
		ctx,
		bson.M{"guildID": guildID, "lists." + args[0]: bson.M{"$exists": true}},
	).Decode(&result)

	if err != nil {
		logger.WithError(err).Error("Error finding list")
		return fmt.Sprintf("Error finding list - %v", err)
	}

	items := result["lists"].(bson.M)[args[0]].(bson.A)

	if len(items) == 0 {
		logger.Debug("Got list")
		return fmt.Sprintf("No items in %v", args[0])
	}
	resultStr := fmt.Sprintf("%v:\n", args[0])

	for _, item := range items {
		resultStr += fmt.Sprintf("\t%v\n", item)
	}

	logger.Debug("Got list")
	return resultStr
}

// getLists returns all of the list names stored by the server.
func (b *Bot) getLists(guildID string) string {
	logger := b.logger.WithFields(log.Fields{
		"guildID": guildID,
		"cmd":     "lists",
	})

	// Get the list
	var result bson.M
	ctx, _ := context.WithTimeout(context.Background(), MongoQueryTimeout*time.Second)
	err := b.collection.FindOne(
		ctx,
		bson.M{"guildID": guildID},
	).Decode(&result)

	if err != nil {
		logger.WithError(err).Error("Error removing lists")
		return fmt.Sprintf("Error finding lists - %v", err)
	}

	lists := result["lists"].(bson.M)
	if len(lists) == 0 {
		logger.Debug("Got lists")
		return "No Lists"
	}
	resultStr := "Lists:\n"

	for list := range lists {
		resultStr += fmt.Sprintf("\t%v\n", list)
	}

	logger.Debug("Got lists")
	return resultStr
}
