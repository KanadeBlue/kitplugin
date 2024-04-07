package main

import (
    "encoding/json"
    "github.com/df-mc/dragonfly/dragonfly/command"
    "github.com/df-mc/dragonfly/dragonfly/player"
    "github.com/df-mc/dragonfly/dragonfly/item"
    "os"
)

type KitItem struct {
    ItemID string `json:"item_id"`
    Amount int    `json:"amount"`
}

type Kit struct {
    Name       string    `json:"name"`
    Items      []KitItem `json:"items"`
    Permission string    `json:"permission"`
}

type KitsConfig struct {
    Kits []Kit `json:"kits"`
}

func main() {
    configFile, err := os.Open("kits.json")
    if err != nil {
        panic("Error opening config file: " + err.Error())
    }
    defer configFile.Close()

    var config KitsConfig
    jsonParser := json.NewDecoder(configFile)
    if err := jsonParser.Decode(&config); err != nil {
        panic("Error parsing config file: " + err.Error())
    }
    command.NewSimpleCommand("kits", func(source command.Source, args []string) {
        if p, ok := source.(*player.Player); ok {
            form := p.NewForm("Kit Selection", "Choose a kit:")
            for _, kit := range config.Kits {
                if p.HasPermission(kit.Permission) {
                    form.AddButton(kit.Name, func() {
                        for _, item := range kit.Items {
                            p.Inventory().Add(itemStackFromID(item.ItemID, item.Amount))
                        }
                        p.SendMessage("You received the " + kit.Name + " kit!")
                    })
                } else {
                    form.AddButton("Locked "+kit.Name, func() {
                        p.SendMessage("This kit is locked.")
                    })
                }
            }
            form.Show(p)
        }
    }).Register(player.PermissionLevel(0))

	player.Interact(func(event player.InteractEvent) {
        if event.Action == player.InteractActionPlaceBlock {
            if chest, ok := event.Target.(*world.Chest); ok {
                if chest.CustomName == "My Kit Chest" {
                    for _, kitItem := range config.KitItems {
                        itemStack, _ := item.NewStackFromNBT(map[string]interface{}{"id": kitItem.ItemID})
                        itemStack.Count = kitItem.Amount
                        chest.Inventory().Add(itemStack)
                    }
                }
            }
        }
    })
}

func itemStackFromID(itemID string, amount int) *item.Stack {
    stack, _ := item.NewStackFromNBT(map[string]interface{}{"id": itemID})
    stack.Count = amount
    return stack
}
