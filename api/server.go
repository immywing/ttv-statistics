package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ttv-statistics/helixclient"
)

var (
	Host string
)

type ttvStatisticsServer struct {
	server http.Server
}

func (s *ttvStatisticsServer) Run() {

	go func() {
		log.Println("\n\t",
			fmt.Sprintf("Serving %s\n\t", apiName),
			fmt.Sprintf("host=%s\n\t", Host),
			fmt.Sprintf("client-id=%s\n\t", helixclient.ClientID),
			fmt.Sprintf("client-secret=%s\n\t", helixclient.ClientSecret),
			fmt.Sprintf("helix-host=%s\n", helixclient.HelixHost),
		)
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("error occurred while serving %s API: %v", apiName, err)
		}
	}()

}

func (s *ttvStatisticsServer) ShutDownServer(ctx context.Context) error {

	log.Printf("shutting down %s API at: %s\n", apiName, s.server.Addr)
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}

func NewTTVStatisticsServer() *ttvStatisticsServer {

	return &ttvStatisticsServer{
		server: http.Server{
			Addr:    Host,
			Handler: wiredMux(),
		},
	}
}

func wiredMux() *http.ServeMux {
	mux := http.NewServeMux()

	for endpoint, handler := range EndpointMapping {
		mux.HandleFunc(endpoint, handler)
	}

	return mux
}
