package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

// Struktura tego z czego słada się user
type User struct {
	Nick       string `json:"Nick"`
	Image      string `json:"Image"`
	UUID       string
	Connection net.Conn
	Permission string `json:"Permission"`
	Rank       string `json:"Rank"`
}

// Z czego skłąda się wiadomość
type Message struct {
	Nick    string `json:"Nick"`
	Message string `json:"Message"`
	Date    string `json:"Date"`
}

// Lista użytkowników ( to jest wiadomośc która jest wysyłana i przez nią widać kto jest online)
type UsersList struct {
	Users []User `json:"Users"`
}

var users []User // Lista użytkowników którą ma serwer

// Połączenie z klientem
func handleConnection(conn net.Conn) {
	defer conn.Close() // Zamykanie połączenia

	buffer := make([]byte, 2000) // Tworzenie bufora na dane

	n, err := conn.Read(buffer) // Odczyt danych
	if err != nil {
		fmt.Println("Błąd odczytu danych:", err)
		return
	}

	user := User{ // Tworzenie nowego użytkownika
		Image:      "Empty",
		UUID:       uuid.New().String(),
		Connection: conn,
	}

	nickJSON := string(buffer[:n])                 // Odczyt nicku
	json.Unmarshal([]byte(nickJSON), &user)        // Parsowanie nicku ( dostajesz nic w wiado )
	fmt.Printf("Nick od klienta: %s\n", user.Nick) // Wyświetlenie nicku

	fmt.Println("Nowy użytkownik:", user.Nick) // Komunikat o nowym użytkowniku

	users = append(users, user) // Dodanie użytkownika do listy

	sendUsersListToAll() // Wysłanie listy użytkowników do wszystkich

	buffer = make([]byte, 2000) // Czyszczenie bufora
	for {                       // Pętla odczytu wiadomości
		n, err := conn.Read(buffer) // Odczyt wiadomości
		if err != nil {
			fmt.Println("Błąd odczytu danych:", err)
			buffer = nil
			break
		}

		message := Message{} // Tworzenie nowej wiadomości

		messageJSON := string(buffer[:n])                               // Odczyt wiadomości
		json.Unmarshal([]byte(messageJSON), &message)                   // Parsowanie wiadomości
		fmt.Printf("Wiadomość od %s: %s\n", user.Nick, message.Message) // Wyświetlenie wiadomości

		message.Nick = user.Nick // Dodanie nicku do wiadomości

		currentTime := time.Now() // Pobranie aktualnej daty

		message.Date = "" + currentTime.Format("2006-01-02 15:04:05") // Dodanie daty do wiadomości

		sendMessageToAll(message) // Wysłanie wiadomości do wszystkich

	}

	removeUser(user)                                 // Usunięcie użytkownika z listy
	fmt.Println("Użytkownik rozłączony:", user.Nick) // Komunikat o rozłączeniu użytkownika

	sendUsersListToAll() // Wysłanie listy użytkowników do wszystkich

}

func main() {

	listener, err := net.Listen("tcp", ":8080") // Uruchomienie serwera TCP
	if err != nil {
		fmt.Println("Błąd uruchamiania serwera TCP:", err)
		return
	}
	defer listener.Close() // Zamknięcie serwera

	fmt.Println("Serwer czatu TCP uruchomiony na porcie 8080") // Komunikat o uruchomieniu serwera

	for { // Pętla oczekiwania na połączenia
		conn, err := listener.Accept() // Akceptowanie połączenia
		if err != nil {
			fmt.Println("Błąd akceptowania połączenia TCP:", err)
			continue
		} else {
			fmt.Println("Nowe połączenie TCP:", conn.RemoteAddr())
		}

		go handleConnection(conn) // Obsługa połączenia
	}
}

func sendUsersListToAll() { // Wysłanie listy użytkowników do wszystkich
	userList := make([]User, len(users)) // Tworzenie listy użytkowników
	for i, user := range users {         // Przepisanie użytkowników do listy
		userList[i] = User{
			Nick:  user.Nick,
			Image: user.Image,
		}
	}

	usersList := UsersList{ // Tworzenie listy użytkowników
		Users: userList,
	}

	usersJSON, err := json.Marshal(usersList) // Parsowanie listy użytkowników
	if err != nil {
		fmt.Println("Error processing data:", err)
		return
	}

	for _, user := range users { // Wysłanie listy użytkowników do wszystkich
		_, err := user.Connection.Write([]byte(usersJSON)) // Wysłanie listy użytkowników
		if err != nil {
			fmt.Println("Error sending data:", err)
		}
	}
}

func sendMessageToAll(message Message) { // Wysłanie wiadomości do wszystkich
	json, err := json.Marshal(message) // Parsowanie wiadomości
	if err != nil {
		fmt.Println("Błąd przetwarzania danych:", err)
		return
	}

	for _, user := range users { // Wysłanie wiadomości do wszystkich
		_, err := user.Connection.Write([]byte(json))
		if err != nil {
			fmt.Println("Błąd wysyłania danych:", err)
		}
	}
}

func removeUser(user User) { // Usunięcie użytkownika z listy
	for i, existingUser := range users { // Szukanie użytkownika
		if existingUser.UUID == user.UUID { // Usunięcie użytkownika
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
}
