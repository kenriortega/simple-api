package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestHandleUsers(t *testing.T) {
	is := is.New(t)

	expectedURL := "https://tSAQgEh.biz/cFnTqrq.php"
	type User struct {
		FirstName   string `json:"first_name,omitempty"`
		LastName    string `json:"last_name,omitempty"`
		Email       string `json:"email,omitempty"`
		Phone       string `json:"phone,omitempty"`
		JobTitle    string `json:"job_title,omitempty"`
		Domain      string `json:"domain,omitempty"`
		URL         string `json:"url,omitempty"`
		PaymentCard string `json:"payment_card,omitempty"`
	}
	var users []User

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(2 * time.Millisecond)
			w.Write([]byte(
				`[
					{
						"_id": "61ec272c19dd549394dcbb36",
						"domain": "PyBuZqZ.biz",
						"email": "ZrBZaeD@qWHMSXN.net",
						"first_name": "Dr. Westley Legros",
						"job_title": "dolor",
						"last_name": "Connelly",
						"payment_card": "MVR 1597662.700000",
						"phone": "104-723-5981",
						"url": "https://tSAQgEh.biz/cFnTqrq.php"
					}
				]`,
			))
		}),
	)
	defer s.Close()

	req, err := http.NewRequest(http.MethodGet, s.URL, nil)
	is.NoErr(err)
	req.Header.Add("content-type", "application/json")
	req = req.WithContext(req.Context())
	resp, err := http.DefaultClient.Do(req)
	is.NoErr(err)
	defer DrainBody(resp.Body)

	data, err := ioutil.ReadAll(resp.Body)
	is.NoErr(err)
	err = json.Unmarshal(data, &users)
	is.NoErr(err)
	is.Equal(expectedURL, users[0].URL)

}

func DrainBody(respBody io.ReadCloser) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if respBody != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		defer respBody.Close()
		_, _ = io.Copy(ioutil.Discard, respBody)
	}
}
