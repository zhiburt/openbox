package httptransport

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/zhiburt/openbox/monitor/services/monitor"
	"github.com/zhiburt/openbox/openbox/monitor/services/monitor/impl"
	"github.com/zhiburt/openbox/openbox/monitor/services/monitor/repositories/mock"
	"github.com/zhiburt/openbox/openbox/monitor/services/monitor/transport/endpoints"
	"github.com/zhiburt/openbox/openbox/monitor/services/qcommunicator"
	qmock "github.com/zhiburt/openbox/openbox/monitor/services/qcommunicator/mock"
)

func TestNewServiceFiles(t *testing.T) {
	cases := []struct {
		qfs                qcommunicator.QFileSystem
		repo               monitor.Repository
		method             string
		url                string
		body               string
		expected           string
		expectedStatusCode int
	}{
		//Create file
		{
			repo:               mock.NewRepository("1111", nil, monitor.File{}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "POST",
			url:                "http://example.com/files",
			body:               `{"owner_id": "123", "name": "some_name"}`,
			expected:           `{"id":"1111"}`,
			expectedStatusCode: 200,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "POST",
			url:                "http://example.com/files",
			body:               `{"owner_id": "123"}`,
			expected:           `{"error":"request error, try to check params"}`,
			expectedStatusCode: 400,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "POST",
			url:                "http://example.com/files",
			body:               ``,
			expected:           `{"error":"EOF"}`,
			expectedStatusCode: 400,
		},
		{
			repo:               mock.NewRepository("", fmt.Errorf("some erorr"), monitor.File{}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "POST",
			url:                "http://example.com/files",
			body:               `{"owner_id": "123", "name": "some_name"}`,
			expected:           `{"error":"error happend in repository"}`,
			expectedStatusCode: 400,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{}),
			qfs:                qmock.NewQcommunicatorMock(nil, fmt.Errorf("some erorr")),
			method:             "POST",
			url:                "http://example.com/files",
			body:               `{"owner_id": "123", "name": "some_name"}`,
			expected:           `{"error":"try to figure out why"}`,
			expectedStatusCode: 400,
		},
		//Lookup file
		{
			repo:               mock.NewRepository("", nil, monitor.File{Name: "1.cpp", OwnerID: "1owner", ServerID: "1"}),
			qfs:                qmock.NewQcommunicatorMock([]byte("hello world"), nil),
			method:             "GET",
			url:                "http://example.com/files/owner/1",
			body:               ``,
			expected:           `{"files":[{"owner_id":"1owner","server_id":"1","name":"1.cpp","body":"aGVsbG8gd29ybGQ=","status":""}]}`,
			expectedStatusCode: 200,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{Name: "1.cpp", OwnerID: "1owner", ServerID: "1"}),
			qfs:                qmock.NewQcommunicatorMock([]byte("hello world"), nil),
			method:             "GET",
			url:                "http://example.com/files/owner/-1",
			body:               ``,
			expected:           `{"files":[{"owner_id":"1owner","server_id":"1","name":"1.cpp","body":"aGVsbG8gd29ybGQ=","status":""}]}`,
			expectedStatusCode: 200,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{Name: "1.cpp", OwnerID: "1owner"}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "GET",
			url:                "http://example.com/files/owner/",
			body:               ``,
			expected:           `404 page not found`,
			expectedStatusCode: 404,
		},
		{
			repo:               mock.NewRepository("", nil, monitor.File{Name: "1.cpp", OwnerID: "1owner", ServerID: "1"}),
			qfs:                qmock.NewQcommunicatorMock(nil, fmt.Errorf("some_error")),
			method:             "GET",
			url:                "http://example.com/files/owner/1",
			body:               ``,
			expected:           `{"error":"error connect with communication with services"}`,
			expectedStatusCode: 400,
		},
		{
			repo:               mock.NewRepository("", fmt.Errorf("some_error"), monitor.File{Name: "1.cpp", OwnerID: "1owner"}),
			qfs:                qmock.NewQcommunicatorMock(nil, nil),
			method:             "GET",
			url:                "http://example.com/files/owner/1",
			body:               ``,
			expected:           `{"error":"error happend in repository"}`,
			expectedStatusCode: 400,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(c.method, c.url, strings.NewReader(c.body))

			monitor := impl.NewService(c.repo, log.NewNopLogger(), c.qfs)
			server := NewService(endpoints.MakeEndpoints(monitor), log.NewNopLogger())

			server.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) != c.expected+"\n" {
				t.Errorf("expected body %s\nbut was %s", c.expected, body)
			}
			if resp.StatusCode != c.expectedStatusCode {
				t.Errorf("expected status code%v but was%v", c.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
