package bot

import (
	"fmt"
	"math/rand"
)

// add appends an item to a specified list. If the list
// does not exist, a new list is created. This function
// returns the text that the bot should send to the server.
func (ds *DiscordServer) add(args []string) string {
	// Check that there are enough arguments given
	if len(args) < 2 {
		return `
		Add Command Error - Incorrect number of arguments
		
		Usage: ~add <list> <item> [items]
		`
	}

	var result string

	// Check the list exists. If not create a new one
	if _, exists := ds.lists[args[0]]; !exists {
		ds.lists[args[0]] = make([]string, 0)
		result += fmt.Sprintf("Created New List: %s\n", args[0])
	}

	// Add the item to the end of the list
	ds.lists[args[0]] = append(ds.lists[args[0]], args[1:]...)
	result += fmt.Sprintf("Added Items: %v", args[1:])

	return result
}

// removeItem takes away one instance of an item from a specified
// list for a server. This function returns the text that the bot
// should send to the server.
func (ds *DiscordServer) removeItem(args []string) string {
	// Check enough arguments for the command were given
	if len(args) < 2 {
		return `
		Remove Command Error - Incorrect number of arguments
		
		Usage: ~remove <list> <item> [items]
		`
	}

	if _, exists := ds.lists[args[0]]; !exists {
		return fmt.Sprintf("Remove Command Error - No list found %s", args[0])
	}

	var result string

	// Get the list which the item should be removed from
	list := ds.lists[args[0]]

	// Loop through the items to be removed
	for _, item := range args[1:] {
		// Set the remove index to something that could
		// never happen
		remove := -1

		// Go through the list and find the index of the item
		// to be removed form the list
		for i := 0; i < len(list); i++ {
			if item == list[i] {
				remove = i
				break
			}
		}

		// Check if the item does not exist in the list
		if remove < 0 || len(list) > remove {
			result += fmt.Sprintf("Could not find Item: %s\n", item)
			continue
		}

		list = RemoveIndex(list, remove)

		result += fmt.Sprintf("Removed Item: %s\n", item)
	}

	// Update the servers version of the list with the
	// list with the removed items
	ds.lists[args[0]] = list

	return result
}

// removeList deletes an entire list from the discord server.
// This function returns the text that the bot should send
// to the server.
func (ds *DiscordServer) removeList(args []string) string {
	// Check enough arguments for the command were given
	if len(args) < 1 {
		return `
		Remove List Command Error - Incorrect number of arguments
		
		Usage: ~removeList <list>
		`
	}

	// Check if the list does not exist
	if _, exists := ds.lists[args[0]]; !exists {
		return fmt.Sprintf("Remove Command Error - No list found %s", args[0])
	}

	delete(ds.lists, args[0])

	return fmt.Sprintf("Removed List: %s", args[0])
}

// pick chooses at random one item from a specified list.
// An optional argument can be provided which will remove
// the chosen item from the list. This function returns
// the text that the bot should send to the server.
func (ds *DiscordServer) pick(args []string) string {
	// Check there are enough arguments
	if len(args) < 1 {
		return "Pick Command Error - Incorrect number of arguments\n\nUsage: ~pick <list> [remove]"
	}

	// Check the specified list exists
	if _, exists := ds.lists[args[0]]; !exists {
		return fmt.Sprintf("Pick Command Error - Could not find list: %s", args[0])
	}

	// Get a random number between 0 and the length of the list.
	// Then get the chosen item from the list.
	chosenIndex := rand.Intn(len(ds.lists[args[0]]))
	chosen := ds.lists[args[0]][chosenIndex]

	// Check the remove argument has been provided
	if len(args) >= 2 {
		// Remove the chosen item from the list
		ds.lists[args[0]] = RemoveIndex(ds.lists[args[0]], chosenIndex)
		return fmt.Sprintf("Random item from list %s chosen and removed: %s", args[0], chosen)
	}

	return fmt.Sprintf("Random item from list %s chosen: %s", args[0], chosen)
}

// getList returns the list of items in the specified list.
func (ds *DiscordServer) getList(args []string) string {
	// Check enough arguments were given
	if len(args) < 1 {
		return `
		List Command Error - Incorrect number of arguments
		
		Usage: ~list <list>
		`
	}

	if _, exists := ds.lists[args[0]]; !exists {
		return fmt.Sprintf("List Command Error - Could not find list: %s", args[0])
	}

	result := fmt.Sprintf("List: %s\n", args[0])

	if len(ds.lists[args[0]]) == 0 {
		result += "\tNo Items"
		return result
	}

	for _, item := range ds.lists[args[0]] {
		result += fmt.Sprintf("\t%s\n", item)
	}

	return result
}

// getLists returns all of the list names stored by the server.
func (ds *DiscordServer) getLists() string {
	result := "Lists:\n"

	if len(ds.lists) == 0 {
		result += "\tNo Lists"
		return result
	}

	for list := range ds.lists {
		result += fmt.Sprintf("\t%s\n", list)
	}

	return result
}
