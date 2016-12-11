package restroute

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"

	"testing"
)

type rmt struct {
	path         string
	shouldMatch  bool
	expectParams map[string]string
}

var routeMatchTestCases = []struct {
	m     Map
	tests []rmt
}{
	{
		Map{
			"/hello/onlin.*": nil,
		},
		[]rmt{
			{"/hello/online", true, nil},
			{"/hello/hajhk", false, nil},
			{"/haih/online", false, nil},
			{"/hello/online_always", true, nil},
		},
	},
	{
		Map{
			`/users/(?P<username>\w+)`: nil,
		},
		[]rmt{
			{"/users/brenk", true, map[string]string{"username": "brenk"}},
			{"/hello/fora", false, map[string]string{"username": "fora"}},
			{"/users/fora", true, map[string]string{"username": "fora"}},
			{"/users/5", true, map[string]string{"username": "5"}},
		},
	},
	{
		Map{
			`/users/(?P<username>\w+)/(?P<field>\w+)`: nil,
		},
		[]rmt{
			{"/users/brenk", false, nil},
			{"/users/5", false, nil},
			{"/users/brenk/date_created", true, map[string]string{"username": "brenk", "field": "date_created"}},
		},
	},
}

func TestRouteMatch(t *testing.T) {
	host := "blah.net"

	for _, c := range routeMatchTestCases {
		routerIf, err := c.m.Compile()
		if err != nil {
			t.Errorf("Failed to compile map: %s", err)
			t.Fail()
			continue
		}

		// Make sure MustCompile works too
		c.m.MustCompile()

		router := routerIf.(*router)

		for _, test := range c.tests {
			url := host + test.path
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Errorf("Could not create request to %q: %s", url, err)
				t.Fail()
				continue
			}

			_, params, ok := router.getRouteFromRequest(req)

			if !ok && test.shouldMatch {
				t.Errorf("Path should match: %q", url)
				t.Fail()
				continue
			} else if ok && !test.shouldMatch {
				t.Errorf("Path should not match: %q", url)
				t.Fail()
				continue
			} else if !ok {
				continue
			}

			if len(params) != len(test.expectParams) {
				t.Errorf("Expected %d params; got %d", len(test.expectParams), len(params))
				t.Fail()
				continue
			}

			for k, v := range test.expectParams {
				if params[k] != v {
					t.Errorf("Wrong param value for %q: wanted %q, got %q", k, v, params[k])
					t.Fail()
					continue
				}
			}
		}
	}
}

type sht struct {
	method       string
	path         string
	expectStatus int
	expectOutput string
}

var serveHTTPTestCases = []struct {
	m     Map
	tests []sht
}{
	{
		Map{
			`/users?/(?P<username>\w+)`: MethodMap{
				"GET": func(r Request) {
					r.W.WriteHeader(200)
					io.WriteString(r.W, fmt.Sprintf("hello %s\n", r.Params["username"]))
				},
				//"POST": nil,
			},
		},
		[]sht{
			{"GET", "/users/beko", 200, "hello beko\n"},
			{"GET", "/user/beko", 200, "hello beko\n"},
			{"GET", "/user/borga", 200, "hello borga\n"},

			{"POST", "/users/beko", 405, ""},

			{"GET", "/usersa/beko", 404, ""},
			{"GET", "/use/beko", 404, ""},
		},
	},
}

type testServer struct {
	port int
	net.Listener
}

func startTestServer(m Map) (*testServer, error) {
	// Get a port that is available for listening on TCP.
	handler, err := m.Compile()
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"), Port: 0,
	})
	if err != nil {
		return nil, err
	}

	port := l.Addr().(*net.TCPAddr).Port

	go http.Serve(l, handler)

	return &testServer{port, l}, nil
}

func TestServeHTTP(t *testing.T) {
	for _, c := range serveHTTPTestCases {
		// Start a server
		s, err := startTestServer(c.m)

		if err != nil {
			t.Errorf("Could not start test server: %s", err)
			t.Fail()
			continue
		}

		defer s.Close()

		client := http.Client{}

		// Make each HTTP request against this server
		for _, test := range c.tests {
			url := fmt.Sprintf("http://localhost:%d%s", s.port, test.path)
			req, err := http.NewRequest(test.method, url, nil)
			if err != nil {
				t.Errorf("Failed to prepare request to '%s %s': %s", test.method, url, err)
				t.Fail()
				continue
			}

			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Failed to make request to '%s %s': %s", test.method, url, err)
				t.Fail()
				continue
			}

			if res.StatusCode != test.expectStatus {
				t.Errorf("Wrong status code for '%s %s': Expected %d, got %d", test.method, url, test.expectStatus, res.StatusCode)
				t.Fail()
				continue
			}
		}
	}
}

func TestBadServer(t *testing.T) {
	// Server with an invalid regexp
	m := Map{
		`\`: MethodMap{},
	}

	_, err := m.Compile()

	if err == nil {
		t.Fatal("Invalid map should fail compile")
	}

	defer func() {
		e := recover()
		if err == nil {
			t.Fatal("Should have panicked")
		}
		if e.(error).Error() != err.Error() {
			t.Fatal("Expected panic from MustCompile() to have the same error as Compile()")
		}
	}()

	m.MustCompile()

	// Should not reach this code path
	t.Error("MustCompile should panic")
	t.FailNow()

}

func fakeHandler(req Request) {
}

var mergeTestCases = []struct {
	maps      []Map
	expectMap Map
}{
	{[]Map{}, Map{}},
	{
		[]Map{
			Map{
				"foo": MethodMap{"GET": fakeHandler},
			},
			Map{
				"bar": MethodMap{"POST": fakeHandler},
			},
		},
		Map{
			"foo": MethodMap{"GET": fakeHandler},
			"bar": MethodMap{"POST": fakeHandler},
		},
	},
}

func TestMerge(t *testing.T) {
	for _, c := range mergeTestCases {
		merged := Merge(c.maps...)

		if len(merged) != len(c.expectMap) {
			t.Fatalf("Merged map did not match expected map")
		}

		for route, mm := range c.expectMap {
			if !reflect.DeepEqual(mm, c.expectMap[route]) {
				t.Fatalf("Merged map did not match expected map")
			}
		}
	}
}
