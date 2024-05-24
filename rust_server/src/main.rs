use std::{io::{Read, Write}, net::{TcpListener, TcpStream}};

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

fn handle_connection(mut stream: TcpStream) {
    // Creating buffer for data reading
    let mut buffer = [0; 1024];
    // Reading data from the stream
    stream.read(&mut buffer).expect("Error reading data");
    // Converting request from bits to utf-8 string and then printing it
    let request = String::from_utf8_lossy(&buffer[..]);
    println!("Received request: {}", request);

    let response = "Hello Client".as_bytes();
    // Sending response to connected peer
    stream.write(response).expect("Cannot read the stream");
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