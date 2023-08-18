package crud

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func chkErr(t *testing.T, resp *http.Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode > 299 {
		t.Fatal(resp.Status)
	}
}

// TestUserCreateAndRead test the User cr of crud
func TestUserCreateAndRead(t *testing.T) {
	usermap := map[string]interface{}{
		"first_name":  "Tom",
		"middle_name": "",
		"last_name":   "Kid",
		"address": map[string]string{
			"line_1":      "111 Easy Street",
			"line_2":      "",
			"city":        "San Francisco",
			"subdivision": "",
			"postal_code": "94530",
		},
	}
	email := "tkid@gmail.com"

	// Clean up previous tests.  No http.Delete ¯\_(ツ)_/¯
	req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/0.9/user/"+email, nil)
	resp, err := http.DefaultClient.Do(req)
	chkErr(t, resp, err)

	data, err := json.Marshal(usermap)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Post("http://localhost:8080/0.9/user", "application/json", bytes.NewBuffer(data))
	if resp.StatusCode != http.StatusNotAcceptable {
		t.Fatalf("Request should of failed due to missing email: %s", err.Error())
	}

	usermap["email"] = email
	data, _ = json.Marshal(usermap)
	resp, err = http.Post("http://localhost:8080/0.9/user", "application/json", bytes.NewBuffer(data))
	chkErr(t, resp, err)

	var cResult map[string]interface{}
	var rResult map[string]interface{}
	var cID, rID string

	json.NewDecoder(resp.Body).Decode(&cResult)
	cID = cResult["id"].(string)

	// This should be a new user
	if !cResult["user_created"].(bool) {
		t.Fatalf("New user wasn't created")
	}

	// One more time for kicks
	resp, err = http.Post("http://localhost:8080/0.9/user", "application/json", bytes.NewBuffer(data))
	chkErr(t, resp, err)
	json.NewDecoder(resp.Body).Decode(&cResult)

	// This should be the same user
	if cResult["user_created"].(bool) {
		t.Fatalf("Ack! We recreated the user?")
	}

	// Try a lookup the user
	resp, err = http.Get("http://localhost:8080/0.9/user/" + cID)
	chkErr(t, resp, err)

	json.NewDecoder(resp.Body).Decode(&rResult)
	rID = rResult["id"].(string)

	if cID != rID {
		t.Fatalf("id query failed: %s != %s", cID, rID)
	}

	resp, err = http.Get("http://localhost:8080/0.9/user/garbage")

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("bad id query didn't 404: %s", err.Error())
	}
}
