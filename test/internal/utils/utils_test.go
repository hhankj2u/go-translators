package utils_test

import (
	"testing"
	"translators/internal/utils"
)

func TestMakeASoup(t *testing.T) {
	validHTML := "<html><body><p>Hello, World!</p></body></html>"
	invalidHTML := "<html><body><p>Hello, World!</body></html>"

	_, err := utils.MakeASoup(validHTML)
	if err != nil {
		t.Fatalf("Expected no error for valid HTML, got %v", err)
	}

	_, err = utils.MakeASoup(invalidHTML)
	if err != nil {
		t.Fatalf("Expected no error for invalid HTML, got %v", err)
	}
}

func TestReplaceAll(t *testing.T) {
	input := "\n            (test)\n      \t                \n      \t                \\''  \xa0 [ ]A1A2B1B2C1C2"
	expected := "(test)  '' []"
	result := utils.ReplaceAll(input)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	input = "no replacements here"
	expected = "no replacements here"
	result = utils.ReplaceAll(input)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestParseResponseURL(t *testing.T) {
	inputURL := "http://example.com/path?query=1"
	expected := "http://example.com/path"
	result := utils.ParseResponseURL(inputURL)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	inputURL = "http://example.com/path"
	expected = "http://example.com/path"
	result = utils.ParseResponseURL(inputURL)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestGetRequestURL(t *testing.T) {
	baseURL := "http://example.com/"
	inputWord := "test word"
	dict := "dict1"
	expected := "http://example.com/test%20word"
	result := utils.GetRequestURL(baseURL, inputWord, dict)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	dict = "dict2"
	expected = "http://example.com/test%20word"
	result = utils.GetRequestURL(baseURL, inputWord, dict)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestGetRequestURLSpellcheck(t *testing.T) {
	baseURL := "http://example.com/"
	inputWord := "test word/test"
	expected := "http://example.com/test+word+test"
	result := utils.GetRequestURLSpellcheck(baseURL, inputWord)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	inputWord = "testword"
	expected = "http://example.com/testword"
	result = utils.GetRequestURLSpellcheck(baseURL, inputWord)
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}
