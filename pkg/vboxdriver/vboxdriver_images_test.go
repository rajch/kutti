package vboxdriver_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/rajch/kutti/pkg/vboxdriver"
)

func TestFetchImageList(t *testing.T) {
	t.Log("Testing FetchImageList...")

	serverMux := http.NewServeMux()
	server := http.Server{Addr: "localhost:8181", Handler: serverMux}
	defer server.Shutdown(context.Background())

	serverMux.HandleFunc(
		"/images.json",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"1.14":{"ImageK8sVersion":"1.14","ImageChecksum":"d071b0f991e4c2ee6b0cd95c77c2ca6336e351717f724236f30b52a31002ff1a","ImageStatus":"Unavailable"},"1.16":{"ImageK8sVersion":"1.16","ImageChecksum":"96460af6aa4fcbf1e9a66e530c044bbfab0d6ce5f453b439287a349eeef1af7d","ImageStatus":"Unavailable"},"1.18":{"ImageK8sVersion":"1.18","ImageChecksum":"dbdfdfaa686143199887d605d6971074886189f68c5fbf081cae874bc9a56da8","ImageStatus":"Unavailable"}}`)
		},
	)

	go func() {
		t.Log("Server starting...")
		err := server.ListenAndServe()
		if err != nil {
			t.Logf("ERROR:%v", err)
		}
		t.Log("Server stopped.")
	}()

	vboxdriver.ImagesSourceURL = "http://localhost:8181/images.json"
	vd, _ := vboxdriver.New()

	t.Log("Fetching image list...")
	err := vd.FetchImageList()
	if err != nil {
		t.Logf("FetchImageList failed with: %v", err)
		t.FailNow()
	}

	t.Log("Fetched image list.")
}
