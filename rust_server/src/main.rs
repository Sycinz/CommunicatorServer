use std::{io::{Read, Write}, net::{TcpListener, TcpStream}};
use uuid::{uuid, Uuid};
use serde_json::{json, Value};

// #[derive(Serialize, Deserialize, Debug)]
struct User {
    nick: String,
    image: String,
    uuid: String,
    connection: TcpStream,
    permission: String,
    rank: String
}
struct Message {
    nick: String,
    message: String,
    date: String
}

struct UsersList {
    users: Vec<String>
}

static mut USERS: Vec<String> = vec![];

fn handle_connection(mut stream: TcpStream) {
    let mut buffer = [0; 2096]; // Creating buffer for data read

    // Reading data from the stream (nickname in this case)
    stream.read(&mut buffer).expect("Error reading data");

    // Converting request from bits to utf-8 string and then printing it
    let request = String::from_utf8_lossy(&buffer[..]);
    println!("Received request: {}", request);

    // Creating new user (need to change it into Map::new()
    let user = User {
        nick: "".to_string(),
        image: "Empty".to_string(),
        uuid: Uuid::new_v4().to_string(),
        connection: stream,
        permission: "".to_string(),
        rank: "".to_string(),
    };

    // Reading nickname from the buffer
    let nick_json = String::from_utf8_lossy(&buffer);
    // JSON unmarshal using serde_json crate
    let nick = serde_json::from_str(nick_json.as_ref())
        .expect("JSON was not properly formatted");
    println!("Nickname of a user: {}", nick_json);

    let mut users: Vec<String> = Vec::new();
    users.push(nick);
}

fn main() {
    let listener = TcpListener::bind("127.0.0.1:3058")
        .expect("Failed to bind IP address");
    println!("Server listening on localhost : 3058");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                std::thread::spawn(|| handle_connection(stream));
            },
            Err(e) => eprintln!("{:?}", e)
        }
    }
}