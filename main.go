package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/CyCoreSystems/dg-whitelist/list"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

//go:embed view/*
var viewfs embed.FS

var dbFile string
var addr string

func init() {
	flag.StringVar(&addr, "addr", ":3000", "address on which to listen")
	flag.StringVar(&dbFile, "f", "/var/lib/dg/dg.db", "database file")
}

func main() {
	flag.Parse()

	db, err := list.Open(dbFile)
	if err != nil {
		log.Println("failed to open database file:", err)
	}
	defer db.Close() // nolint: errcheck

	engine := html.NewFileSystem(http.FS(viewfs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		listName := c.Query("list", list.ListWhite)

		items, err := db.Get(listName)
		if err != nil {
			return c.SendString(fmt.Sprintf("failed to get list %q: %s", listName, err.Error()))
		}

		return c.Render("view/index", fiber.Map{
			"List": listName,
			"Items": items,
		})
	})

	app.Get("/:list", func(c *fiber.Ctx) error {
		list, err := db.Get(c.Query("list", list.ListWhite))
		if err != nil {
			return err
		}

		return c.JSON(list)
	})

	app.Post("/", func (c *fiber.Ctx) error {
		i := new(list.Item)

		if err := c.BodyParser(i); err != nil {
			log.Println("failed to parse body:", err)
			return c.SendStatus(http.StatusBadRequest)
		}

		switch i.List {
		case list.ListWhite:
		case list.ListGrey:
		case list.ListBlack:
		default:
			return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("invalid list %q", i.List))
		}

		// op is sent with the web client
		switch c.FormValue("op", "") {
		case "add":
			if err := db.Add(i.List, i.Address); err != nil {
				log.Printf("failed to add %q to %q: %s", i.Address, i.List, err.Error())
				return c.SendStatus(http.StatusInternalServerError)
			}
			return c.Redirect(fmt.Sprintf("/?list=%s", i.List), http.StatusSeeOther)
		case "delete":
			if err := db.Remove(i.List, i.Address); err != nil {
				log.Printf("failed to remove %q to %q: %s", i.Address, i.List, err.Error())
				return c.SendStatus(http.StatusInternalServerError)
			}
			return c.Redirect(fmt.Sprintf("/?list=%s", i.List), http.StatusSeeOther)
		default:
			// CLI API client
		}

		if err := db.Add(i.List, i.Address); err != nil {
			log.Printf("failed to add %q to %q: %s", i.Address, i.List, err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusNoContent)
	})

	app.Delete("/", func (c *fiber.Ctx) error {
		i := new(list.Item)

		if err := c.BodyParser(i); err != nil {
			log.Println("failed to parse body", err)
			return c.SendStatus(http.StatusBadRequest)
		}
		
		if err := db.Remove(i.List, i.Address); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.SendStatus(http.StatusNoContent)
	})

	log.Fatalln(app.Listen(addr))
}
