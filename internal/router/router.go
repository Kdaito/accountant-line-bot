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

func (r *Router) Set(isSkipGpt bool) {
	// port setting
	if r.Port == "" {
		r.Port = "2001"
	}

	ctx := context.Background()

	pkgServices := setup.NewPkgServices(ctx, isSkipGpt)

	// DI
	drivePkg := gcp.NewDrive(pkgServices.Drive())
	sheetPkg := gcp.NewSheet(pkgServices.Sheet())
	chatAIPkg := gpt.NewChatAI(pkgServices.Gpt())
	messagePkg := line.NewMessage(pkgServices.LineBot())

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
