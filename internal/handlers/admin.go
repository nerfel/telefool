package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"telefool/pkg/di"
)

const createUserRegexpString = `^create-user\s+(@\w+)$`

func CreateUserRoute(ctx *di.UpdateContext) bool {
	createUserRegexp := regexp.MustCompile(createUserRegexpString)
	textMessage := strings.TrimSpace(ctx.Update.Message.Text)

	if !createUserRegexp.MatchString(textMessage) {
		return false
	}

	matches := createUserRegexp.FindStringSubmatch(textMessage)
	if len(matches) < 2 {
		return false
	}
	ctx.RoutePayload = matches[1]

	return true
}

func CreateUserHandler(ctx *di.UpdateContext) {
	userName := ctx.RoutePayload
	fmt.Println("route payload: ", userName)

}
