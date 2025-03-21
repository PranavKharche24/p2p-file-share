package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const PORT = ":8080"
const BUFFER_SIZE = 4096

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Select an option:")
	fmt.Println("[1] Act as Receiver (Server)")
	fmt.Println("[2] Act as Sender (Client)")
	fmt.Print("Enter choice (1 or 2): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "1" {
		startServer()
	} else if choice == "2" {
		startClient()
	} else {
		fmt.Println("Invalid choice. Exiting.")
	}
}

// Receiver (Server)
func startServer() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Receiver ready, waiting for sender connection on port", PORT)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connection established from", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	fileName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}
	fileName = strings.TrimSpace(fileName)
	filePath := filepath.Join("files", fileName)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	bytesWritten, err := io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Printf("File %s received successfully (%d bytes).\n", fileName, bytesWritten)
}

// Sender (Client)
func startClient() {
	consoleReader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter receiver IP address: ")
	receiverIP, _ := consoleReader.ReadString('\n')
	receiverIP = strings.TrimSpace(receiverIP)

	conn, err := net.Dial("tcp", receiverIP+PORT)
	if err != nil {
		fmt.Println("Error connecting to receiver:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Enter path of file to send: ")
	filePath, _ := consoleReader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	fmt.Fprintf(conn, fileName+"\n")

	bytesSent, err := io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}

	fmt.Printf("File %s sent successfully (%d bytes).\n", fileName, bytesSent)
}
