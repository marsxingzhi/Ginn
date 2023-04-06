package ginn

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestParsePath(t *testing.T) {
	convey.Convey("simple test", t, func() {

		var path = "/v1/user/login"

		got := parsePath(path)

		want := []string{"v1", "user", "login"}

		// ShouldResemble 用于数组、切片、map和结构体相等
		convey.So(got, convey.ShouldResemble, want)

	})
}

func TestGetRouter(t *testing.T) {
	r := newTestRouter()

	convey.Convey("test simple", t, func() {
		node := r.getRoute("GET", "/user/login")
		convey.So(node, convey.ShouldNotBeNil)
	})

	convey.Convey("method failed", t, func() {
		node := r.getRoute("GET", "/user/update")
		convey.So(node, convey.ShouldBeNil)
	})

	convey.Convey("method success", t, func() {
		node := r.getRoute("POST", "/user/update")
		convey.So(node, convey.ShouldNotBeNil)
	})

}

func newTestRouter() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/user/login", nil)
	r.addRouter("GET", "/user/register", nil)
	r.addRouter("POST", "/user/update", nil)

	return r
}
