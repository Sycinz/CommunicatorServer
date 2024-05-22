use std::{io::{Read, Write}, net::{TcpListener, TcpStream}};

fn handle_connection(mut stream: TcpStream) {
    let mut buffer = [0; 1024];

    stream.read(&mut buffer).expect("Error reading data");

    let request = String::from_utf8_lossy(&buffer[..]);
    println!("Received request: {}", request);

    let response = "Hello Client".as_bytes();

    stream.write(response).expect("Mazno ni");
    stream.write(response).expect("Problema");
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