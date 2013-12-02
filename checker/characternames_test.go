package checker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"testing"
)

func TestIsCharacterReferenceName(t *testing.T) {

	names := []string{
		"amp",
		"AMP",
		"And",
		"abreve",
		"ZeroWidthSpace",
		"lt", // smallest known
		"CounterClockwiseContourIntegral", // longest known
	}
	casesShouldBeTrue(t, names, IsCharacterReferenceName,
		"Expected %#v to be a character name, but got false")

	notNames := []string{
		"Amp",
		"Daniel",
		"Tuesday",
		"",
		string(UnicodePOI),
	}
	casesShouldBeFalse(t, notNames, IsCharacterReferenceName,
		"Expected %#v to NOT be a character name, but got true")
}

func ExampleIsCharacterReferenceName() {
	fmt.Println(IsCharacterReferenceName("amp"))
	fmt.Println(IsCharacterReferenceName("AMP"))
	fmt.Println(IsCharacterReferenceName("Amp")) // Names are case sensitive.
	fmt.Println(IsCharacterReferenceName("&amp;"))
	// Output:
	// true
	// true
	// false
	// false
}

// BenchmarkIsCharacterReferenceNameTrue	 5000000	     351.0 ns/op	     0 B/op	       0 allocs/op # search sorted array
// BenchmarkIsCharacterReferenceNameTrue	50000000	      50.4 ns/op	     0 B/op	       0 allocs/op # map[string]bool
// BenchmarkIsCharacterReferenceNameTrue    50000000          22.9 ns/op         0 B/op        0 allocs/op # Go 1.2
func BenchmarkIsCharacterReferenceNameTrue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsCharacterReferenceName("nbsp")
	}
}

// BenchmarkIsCharacterReferenceNameFalse	 5000000	      347.0 ns/op	       0 B/op	       0 allocs/op # search sorted array
// BenchmarkIsCharacterReferenceNameFalse	20000000	       79.6 ns/op	       0 B/op	       0 allocs/op # map[string]bool
// BenchmarkIsCharacterReferenceNameFalse  50000000            38.2 ns/op          0 B/op          0 allocs/op # Go 1.2
func BenchmarkIsCharacterReferenceNameFalse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsCharacterReferenceName("tuesday")
	}
}

func TestIsCharacterReference(t *testing.T) {

	names := []string{
		"&amp;",
		"&AMP;",
		"&And;",
		"&abreve;",
		"&ZeroWidthSpace;",
		"&lt;", // smallest known
		"&CounterClockwiseContourIntegral;", // longest known
	}
	casesShouldBeTrue(t, names, IsCharacterReference,
		"Expected %#v to be a character reference, but got false")

	notNames := []string{
		"&Amp;",
		"&Daniel;",
		"&Tuesday;",
		"amp",
		"&\u2318;",
		"",
	}
	casesShouldBeFalse(t, notNames, IsCharacterReference,
		"Expected %#v to NOT be a character reference, but got true")
}

func ExampleIsCharacterReference() {
	fmt.Println(IsCharacterReference("&amp;"))
	fmt.Println(IsCharacterReference("&AMP;"))
	fmt.Println(IsCharacterReference("amp"))   // Needs & and ;
	fmt.Println(IsCharacterReference("&Amp;")) // Case sensitive
	// Output:
	// true
	// true
	// false
	// false
}

func casesShouldBeTrue(t *testing.T, cases []string, test func(string) bool, pattern string) {
	for _, arg := range cases {
		if test(arg) != true {
			t.Errorf(pattern, arg)
		}
	}
}

func casesShouldBeFalse(t *testing.T, cases []string, test func(string) bool, pattern string) {
	for _, arg := range cases {
		if test(arg) != false {
			t.Errorf(pattern, arg)
		}
	}
}

// TestDownloadEntitiesJson checks the characterReferenceNames map against the
// list of entities from the WHATWG to make sure they are in sync.
//
// This test should be skipped by default and only enabled and run deliberately.
//
func xxTestDownloadEntitiesJson(t *testing.T) {

	t.Skip("Skipping TestGetData because it takes a long time.")

	if testing.Short() {
		t.Skip("Skipping TestGetData in short mode.")
	}

	t.Log("Checking characterReferenceNames against the WHATWG list.")

	entitiesUrl := `http://www.whatwg.org/specs/web-apps/current-work/multipage/entities.json`
	t.Log("URL", entitiesUrl)

	// Download entities.json

	resp, err := http.Get(entitiesUrl)
	t.Logf("HTTP Status: %d", resp.StatusCode)
	t.Logf("ETag: %s", resp.Header.Get("ETag"))
	// ETag: 239e9-4ea257efc5c00
	//
	// TODO: Do a HEAD request and only download and parse the JSON file if the
	// ETag is different from the last one?
	//
	if err != nil {
		t.Fatalf("Unable to fetch JSON file: %s", err)
	}
	defer resp.Body.Close()

	// Extract JSON from HTTP response.

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read in HTTP response: %s", err)
	}
	t.Logf("JSON text is %d bytes.", len(httpBody))

	// Parse JSON.

	var v map[string]interface{}
	err = json.Unmarshal(httpBody, &v)
	if err != nil {
		t.Fatalf("Unable to unmarshal JSON data: %s", err)
	}

	// Convert unwieldy JSON object into simple slice

	expectedNames := keys(v)
	for i, name := range expectedNames {
		expectedNames[i] = strings.Trim(name, `&;`)
	}
	t.Logf("Found %d character names.", len(expectedNames))

	// Check each entity against our list. Make sure it exists in the list.

	for i, wanted := range expectedNames {
		if characterReferenceNames[wanted] {
			// Name exists. Delete it so we can see if there are any left-overs.
			delete(characterReferenceNames, wanted)
		} else {
			// There are some duplicate entries in entities.json.
			if i > 0 && expectedNames[i-1] == wanted {
				t.Logf("Name %s is duplicate.", wanted)
			} else {
				t.Errorf("Expected %s to be a character reference name, but got false instead.", wanted)
			}
		}
	}

	// Make sure there are no remaining entries in characterReferenceNames that
	// should be deleted.

	if len(characterReferenceNames) > 0 {
		t.Error("The following values are listed as character reference names, but they were not found in the entities.json list:",
			characterReferenceNames)
	}

	t.Log("Done.")
}

func keys(m map[string]interface{}) []string {
	str := make([]string, len(m))

	i := 0
	for key, _ := range m {
		str[i] = key
		i++
	}

	sort.Strings(str)
	return str
}
