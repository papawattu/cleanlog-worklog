package main

import (
	"net"
)

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

// func Test_get_ok(t *testing.T) {
// 	workLogRepo := NewWorkLogRepository()
// 	workService := NewWorkService(workLogRepo)

// 	server := &http.Server{
// 		Addr: ":0",
// 	}

// 	// convert int to string
// 	if err := startWebServer(server, workService); err != nil {
// 		t.Fatal(err)
// 	}

// 	log.Printf("Starting Work Log server on port %s\n", server.Addr)

// 	r, err := http.Get("http://localhost" + server.Addr + "/api/worklog")

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if r.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code %v, got %v", http.StatusOK, r.StatusCode)
// 	}

// 	t.Log("Test passed")

// }
