package telegram

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_GetUpdates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
		{
		"ok": true,
		"result": [
			{
				"update_id": 123,
				"message": {
					"text": "test",
					"chat": { 
						"id": 111
					}
				}
			}
		]
		}
		`))
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}

	client := NewClient(u.Scheme, u.Host, "test-token")
	updates, err := client.GetUpdates(0, 1)
	if err != nil {
		t.Fatalf("GetUpdates failed: %v", err)
	}

	if len(updates) != 1 || updates[0].Message.Text != "test" {
		t.Errorf("unexpected update: %+v", updates)
	}
}

func TestClient_SendMessage(t *testing.T) {
	var receivedQuery url.Values
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		receivedPath = r.URL.Path

		_, _ = w.Write([]byte(`
		{
 		"ok": true,
 		"result": {
    		"text": "your test",
    		"chat": {
      			"id": 101
    		}
  		}
		}
		`))
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}

	client := NewClient(u.Scheme, u.Host, "test-token")
	err = client.SendMessage(101, "your test")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if receivedPath != "/bottest-token/sendMessage" {
		t.Errorf("unexpected path: got %s", receivedPath)
	}

	if receivedQuery.Get("chat_id") != "101" || receivedQuery.Get("text") != "your test" {
		t.Errorf("unexpected query: %v", receivedQuery)
	}
}
