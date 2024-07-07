use std::{io::{Read, Write}, net::{TcpListener, TcpStream}};
use std::fmt::format;
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
    let mut buffer = Vec::new(); // Creating buffer for data read

    // Reading data from the stream (nickname in this case)
    stream.read_to_end(&mut buffer).expect("Error reading data");

    // Converting request from bits to utf-8 string and then printing it
    let request = String::from_utf8_lossy(&buffer[..]);
    println!("Received request: {}", request);

    let uuid = Uuid::new_v4().to_string();
    // Creating new user (need to change it into Map::new()
    let mut user = User {
        nick: serde_json::from_str("").unwrap(),
        image: serde_json::from_str("Empty").unwrap(),
        uuid: serde_json::from_str(uuid.as_str()).unwrap(),
        connection: stream,
        permission: serde_json::from_str("").unwrap(),
        rank: serde_json::from_str("").unwrap(),
    };

    // Reading nickname from the buffer
    let nick_json = String::from_utf8_lossy(&buffer);

    // JSON unmarshal using serde_json crate
    let nick: String = serde_json::from_str(nick_json.as_ref())
        .expect("JSON was not properly formatted");

    println!("Nickname of a user: {}", nick_json);

    let user_struct_list: Vec<User> = vec![user];

    // send_userlist_to_all(user_struct_list);
}

fn send_userlist_to_all(users_list: Vec<User>) -> String {
    let json = String::new();

    let _ = users_list.iter().map(|user| {
        
    });

    r"Mazno ni".to_string()
}

fn main() {
    let listener = TcpListener::bind("127.0.0.1:8080")
        .expect("Failed to bind IP address");
    println!("Server listening on localhost : 8080");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                std::thread::spawn(|| handle_connection(stream));
            },
            Err(e) => eprintln!("{:?}", e)
        }
    }
}