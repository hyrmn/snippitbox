package main

import "github.com/hyrmn/snippetbox/pkg/models"

type App struct {
	Port      string
	HTMLDir   string
	StaticDir string
	Store     *models.Store
}
