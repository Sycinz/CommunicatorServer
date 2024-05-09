use std::net::TcpListener;

fn main() {
    let listener = TcpListener::bind("127.0.0.1:7878").unwrap();

    handle_connection(listener);
}

fn handle_connection(listener: TcpListener) {

    for stream in listener.incoming() {
        let stream = stream.unwrap();

        println!("Connection established!");
    }
}