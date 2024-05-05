package main

import (
    "encoding/json"
    "fmt"
    "net"
    "time"
    "github.com/google/uuid"
	"gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type User struct {
    Nick       string `json:"Nick"`
    Image      string `json:"Image"`
    UUID       string
    Connection net.Conn
	Permission string `json:"Permission"`
	Rank 	   string `json:"Rank"`
}

type Message struct {
    Nick    string `json:"Nick"`
    Message string `json:"Message"`
    Date    string `json:"Date"`
}

type UsersList struct {
    Users []User `json:"Users"`
}

type Type int

const (
	Chat Type = 1
	Voice Type = 2
)

type Channel struct {
	Users []User
	MaximumUsers int
	ChannelName string
	ChannelDescription string
	ChannelPassword string
	ChannelType Type
	CreationTime time.Time
}

var users []User

func handleConnection(conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 2000)

    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Błąd odczytu danych:", err)
        return
    }

    user := User{
        Image:      "Empty",
        UUID:       uuid.New().String(),
        Connection: conn,
    }

    nickJSON := string(buffer[:n])
    json.Unmarshal([]byte(nickJSON), &user)
    fmt.Printf("Nick od klienta: %s\n", user.Nick)

    fmt.Println("Nowy użytkownik:", user.Nick)

    users = append(users, user)

    sendUsersListToAll()

    buffer = make([]byte, 2000)
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Błąd odczytu danych:", err)
            buffer = nil
            break
        }

        message := Message{}

        messageJSON := string(buffer[:n])
        json.Unmarshal([]byte(messageJSON), &message)
        fmt.Printf("Wiadomość od %s: %s\n", user.Nick, message.Message)

        message.Nick = user.Nick

        currentTime := time.Now()

        message.Date = "" + currentTime.Format("2006-01-02 15:04:05")

        sendMessageToAll(message)

    }

    removeUser(user)
    fmt.Println("Użytkownik rozłączony:", user.Nick)

    sendUsersListToAll()

    message := Message{}
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
	// Check if server.db file exists, if not create it
	db, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
  	if err != nil {
    	panic("failed to connect database")
  	}

	// Create tables if they don't exist
	db.AutoMigrate(&Channel{})

	// TCP server
	go func() {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			fmt.Println("Błąd uruchamiania serwera TCP:", err)
			return
		}
		defer listener.Close()

		fmt.Println("Serwer czatu TCP uruchomiony na porcie 8080")

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Błąd akceptowania połączenia TCP:", err)
				continue
			} else {
				fmt.Println("Nowe połączenie TCP:", conn.RemoteAddr())
			}

			go handleConnection(conn)
		}
	}()

	// UDP server
	go func() {
		udpAddr, err := net.ResolveUDPAddr("udp", ":8081")
		if err != nil {
			fmt.Println("Błąd uruchamiania serwera UDP:", err)
			return
		}

		udpConn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Println("Błąd uruchamiania serwera UDP:", err)
			return
		}
		defer udpConn.Close()

		fmt.Println("Serwer przesyłania głosu UDP uruchomiony na porcie 8081")

		buffer := make([]byte, 1024)

		for {
			n, addr, err := udpConn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Błąd odczytu danych UDP:", err)
				continue
			} else {
				fmt.Println("Nowe połączenie UDP:", addr)
			}

			// Handle voice data here
			// ...

			// Example: Echo back the received data
			_, err = udpConn.WriteToUDP(buffer[:n], addr)
			if err != nil {
				fmt.Println("Błąd wysyłania danych UDP:", err)
			}
		}
	}()

	// Wait indefinitely
	select {}
}

func sendUsersListToAll() {
	userList := make([]User, len(users))
	for i, user := range users {
		userList[i] = User{
			Nick:  user.Nick,
			Image: user.Image,
		}
	}

	usersList := UsersList{
		Users: userList,
	}

	usersJSON, err := json.Marshal(usersList)
	if err != nil {
		fmt.Println("Error processing data:", err)
		return
	}

	for _, user := range users {
		_, err := user.Connection.Write([]byte(usersJSON))
		if err != nil {
			fmt.Println("Error sending data:", err)
		}
	}
}

func sendMessageToAll(message Message) {
    json, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Błąd przetwarzania danych:", err)
        return
    }

    for _, user := range users {
        _, err := user.Connection.Write([]byte(json))
        if err != nil {
            fmt.Println("Błąd wysyłania danych:", err)
        }
    }
}

func removeUser(user User) {
    for i, existingUser := range users {
        if existingUser.UUID == user.UUID {
            users = append(users[:i], users[i+1:]...)
            break
        }
    }
}