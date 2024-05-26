package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/dreamilk/rpc_server/log"
)

func healthCheck() {
	ctx := context.Background()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		msg := fmt.Sprintf("time:%s status: ok", time.Now().Format("2006-01-02 15:04:05"))

		fmt.Fprint(w, msg)
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Error(ctx, "", zap.Error(err))
	}
}
