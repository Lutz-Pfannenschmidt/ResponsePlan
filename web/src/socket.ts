export let socket: WebSocket;
export const dev = process.env.NODE_ENV === 'development';

console.log("DEV:", dev);
if (dev) {
    socket = new WebSocket("ws://127.0.0.1:1337/ws");
} else {
    socket = new WebSocket(`ws://${location.host}/ws`);
}


socket.onopen = () => {
    console.log("Successfully Connected");
};

socket.onmessage = event => {
    console.log(event);
}

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};