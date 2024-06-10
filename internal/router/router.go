package router

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/pkg/gcp"
	"github.com/Kdaito/accountant-line-bot/internal/pkg/gpt"
	"github.com/Kdaito/accountant-line-bot/internal/pkg/line"
	"github.com/Kdaito/accountant-line-bot/internal/pkg/setup"
	"github.com/Kdaito/accountant-line-bot/internal/service"
)

type Router struct {
	Port string
}

func (r *Router) Set() {
	// port setting
	if r.Port == "" {
		r.Port = "2001"
	}

	ctx := context.Background()

	pkgServices := setup.NewPkgServices(ctx)

	// DI
	drivePkg := gcp.NewDrive(pkgServices.GetDrive())
	sheetPkg := gcp.NewSheet(pkgServices.GetSheet())
	chatAIPkg := gpt.NewChatAI(pkgServices.GetGpt())
	messagePkg := line.NewMessage(pkgServices.GetLineBot())

	callbackService := service.NewCallbackService(drivePkg, messagePkg, sheetPkg, chatAIPkg)
	healthService := service.NewHealthService()

	// set routing
	http.HandleFunc("/callback", callbackService.Callback)
	http.HandleFunc("/health-check", healthService.HealthCheck)
}

func (r *Router) Run() {
	fmt.Println("http://localhost:" + r.Port + "/")
	if err := http.ListenAndServe(":"+r.Port, nil); err != nil {
		log.Fatal(err)
	}
}
