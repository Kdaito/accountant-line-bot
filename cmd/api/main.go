// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/handler"
	"github.com/Kdaito/accountant-line-bot/pkg"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func main() {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	port := os.Getenv("PORT")

	bot, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	processer := pkg.NewProcesser(channelSecret, bot)
	handler := handler.NewHandler(processer)

	http.HandleFunc("/callback", handler.HandleCallback)

	if port == "" {
		port = "2001"
	}

	fmt.Println("http://localhost:" + port + "/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}