package ginn

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Panicf("%s\n", trace(msg))
				ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		ctx.Next()
	}
}

func trace(msg string) string {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)

	var str strings.Builder

	str.WriteString(msg + "\nTraceback:")

	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
