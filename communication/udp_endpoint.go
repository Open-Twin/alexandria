package communication

import (
	"fmt"
	"runtime"

	"net"
	"os"
)

/**
Starts the UDP endpoint
 */
func StartUDP() {
	addr := net.UDPAddr{
		IP: []byte{
			0, 0, 0, 0,
		},
		Port: 53,
		Zone: "",
	}
	// UDP Listener
	conn, err := net.ListenUDP("udp", &addr)
	//checks if errors are thrown
	checkError(err)

	quit := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		//starts a new thread that reads from the connection
		go listen(conn, quit)
	}
	<-quit // hang until an error
}
/**
reads from the udp connection
 */
func listen(connection *net.UDPConn, quit chan struct{}) {
	buffer := make([]byte, 1024)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for err == nil {
		n, remoteAddr, err = connection.ReadFromUDP(buffer)
		go handleData(n, buffer, remoteAddr, connection)
		// you might copy out the contents of the packet here, to
		// `var r myapp.Request`, say, and `go handleRequest(r)` (or
		// send it down a channel) to free up the listening
		// goroutine. you do *need* to copy then, though,
		// because you've only made one buffer per listen().
		fmt.Println("from", remoteAddr, "-", buffer[:n])
	}
	fmt.Println("listener failed - ", err)
	quit <- struct{}{}
}
/**
Handles the data from the udp connection
 */
func handleData(n int, buffer []byte, addr* net.UDPAddr, conn *net.UDPConn){

	fmt.Printf("\n--------------\n")
	fmt.Printf("packet-received: bytes=%d from=%s\n",
		n, addr.String())
	fmt.Println("from", addr, "-", buffer[:n])
	fmt.Printf("\n--------------\n")
	//Writes back to the client
	_, err2 := conn.WriteTo(buffer[0:n],addr)
	if err2 != nil {
		return
	}
}
/**
Checks if errors are thrown. If yes, it prints the error and exits the program
*/
func checkError(err error){

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//Exits the program
		os.Exit(1)
	}
}
