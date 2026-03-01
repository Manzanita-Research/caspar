package ghost

import "testing"

func TestIsID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"5ddc9141c35e7700383b2937", true},
		{"abcdef0123456789abcdef01", true},
		{"my-post-slug", false},
		{"hello-world", false},
		{"", false},
		{"5ddc9141c35e7700383b293", false},  // too short
		{"5ddc9141c35e7700383b29377", false}, // too long
		{"5ddc9141c35e7700383b293g", false},  // invalid hex char
	}

	for _, tt := range tests {
		got := IsID(tt.input)
		if got != tt.want {
			t.Errorf("IsID(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestListParams_ToValues(t *testing.T) {
	params := ListParams{
		Limit:   10,
		Page:    2,
		Filter:  "status:published",
		Order:   "published_at desc",
		Fields:  "id,title,slug",
		Include: "tags,authors",
	}

	v := params.toValues()

	if v.Get("limit") != "10" {
		t.Errorf("expected limit=10, got %s", v.Get("limit"))
	}
	if v.Get("page") != "2" {
		t.Errorf("expected page=2, got %s", v.Get("page"))
	}
	if v.Get("filter") != "status:published" {
		t.Errorf("expected filter=status:published, got %s", v.Get("filter"))
	}
	if v.Get("order") != "published_at desc" {
		t.Errorf("expected order=published_at desc, got %s", v.Get("order"))
	}
	if v.Get("fields") != "id,title,slug" {
		t.Errorf("expected fields=id,title,slug, got %s", v.Get("fields"))
	}
	if v.Get("include") != "tags,authors" {
		t.Errorf("expected include=tags,authors, got %s", v.Get("include"))
	}
}

func TestListParams_EmptyValues(t *testing.T) {
	params := ListParams{}
	v := params.toValues()

	if len(v) != 0 {
		t.Errorf("expected empty values for zero params, got %v", v)
	}
}
