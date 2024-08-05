package anthropic //nolint:testpackage // testing private field

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestRequestBuilderReturnsRequest(t *testing.T) {
	b := NewRequestBuilder()
	var (
		ctx         = context.Background()
		method      = http.MethodPost
		url         = "/foo"
		request     = map[string]string{"foo": "bar"}
		reqBytes, _ = json.Marshal(request)
		want, _     = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBytes))
	)
	got, _ := b.Build(ctx, method, url, request, nil)
	if !reflect.DeepEqual(got.Body, want.Body) ||
		!reflect.DeepEqual(got.URL, want.URL) ||
		!reflect.DeepEqual(got.Method, want.Method) {
		t.Errorf("Build() got = %v, want %v", got, want)
	}
}

func TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil(t *testing.T) {
	var (
		ctx     = context.Background()
		method  = http.MethodGet
		url     = "/foo"
		want, _ = http.NewRequestWithContext(ctx, method, url, nil)
	)
	b := NewRequestBuilder()
	got, _ := b.Build(ctx, method, url, nil, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() got = %v, want %v", got, want)
	}
}
