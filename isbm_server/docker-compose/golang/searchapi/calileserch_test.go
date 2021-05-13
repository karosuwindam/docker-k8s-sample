package searchapi

import "testing"

func TestCalilSearch(t *testing.T) {
	isbn_sm := []string{"9784065197653", "4560158370760", "9784758068833", "4910163990536"}
	for _, isbn := range isbn_sm {
		output := GetPageCalilURL(isbn)
		if output.Title != "" {
			t.Log(output)
		}
	}
}
