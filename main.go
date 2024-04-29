package main

import (
	"fmt"
	"net"
	"time"
	// "strings"
	"github.com/google/uuid"
	"encoding/json"
)

type User struct {
	Nick    string `json:"Nick"`
	Image   string `json:"Image"`
	UUID    string
	Connection net.Conn
	// IPv4    string
}

type Message struct {
	Nick    string `json:"Nick"`
	Message string `json:"Message"`
	Date	string `json:"Date"`
}

var users []User

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2000)

	// Odbierz JSON od klienta który ma tylko wartość "Nick"
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Błąd odczytu danych:", err)
		return
	}

	user := User{
		Image: "Empty",
		UUID: uuid.New().String(),
		Connection: conn,
	}

	// Przetwarzanie JSON
	nickJSON := string(buffer[:n])
	json.Unmarshal([]byte(nickJSON), &user)
	fmt.Printf("Nick od klienta: %s\n", user.Nick)

	fmt.Println("Nowy użytkownik:", user.Nick)

	// Dodawanie użytkownika do tablicy
	users = append(users, user)

	var usersList string = ""

	for x, otherUser := range users {
		usersList += "{\"Nick\": \"" + otherUser.Nick + "\", \"Image\": \"" + otherUser.Image + "\"}"
		if x < len(users) - 1 {
			usersList += ", "
		}
	}

	var usersJSON string = "{ \"Users\": [" + usersList + "] }"

	_, err = conn.Write([]byte(usersJSON))

	// Odbieranie wiadomości od klienta
	buffer = make([]byte, 2000)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Błąd odczytu danych:", err)
			buffer = nil
			break
		}

		message := Message{};

		// Przetwarzanie wiadomości
		messageJSON := string(buffer[:n])
		json.Unmarshal([]byte(messageJSON), &message)
		fmt.Printf("Wiadomość od %s: %s\n", user.Nick, message.Message)

		message.Nick = user.Nick

		currentTime := time.Now()

		message.Date = "" + currentTime.Format("2006-01-02 15:04:05")

		// Wysyłanie wiadomości do innych użytkowników
		for _, _ = range users {
			// if otherUser.UUID != user.UUID {
				json, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Błąd przetwarzania danych:", err)
					buffer = nil
					break
				}

				_, err = conn.Write([]byte(json))
				if err != nil {
					fmt.Println("Błąd wysyłania danych:", err)
					buffer = nil
					break
				}
			// }
		}
	}

	// Usuwanie użytkownika po rozłączeniu
	for i, existingUser := range users {
		if existingUser.UUID == user.UUID {
			// Usuń użytkownika z listy
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	fmt.Println("Użytkownik rozłączony:", user.Nick)

	message := Message{};
	message.Nick = "Server"
	message.Message = user.Nick + " rozłączył się"
	currentTime := time.Now()
	message.Date = "" + currentTime.Format("2006-01-02 15:04:05")

	json, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Błąd przetwarzania danych:", err)
	}

	_, err = conn.Write([]byte(json))
	if err != nil {
		fmt.Println("Błąd wysyłania danych:", err)
	}

	conn.Close()

}

func main() {
	// Uruchamianie serwera na porcie 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Błąd uruchamiania serwera:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serwer czatu TCP uruchomiony na porcie 8080")

	// Akceptowanie połączeń od klientów
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Błąd akceptowania połączenia:", err)
			continue
		} else {
			fmt.Println("Nowe połączenie:", conn.RemoteAddr())
		}

		// Obsługa połączenia w nowym wątku
		go handleConnection(conn)
	}
}