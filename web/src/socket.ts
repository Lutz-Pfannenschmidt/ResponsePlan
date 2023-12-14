export let socket = new WebSocket("ws://127.0.0.1:1337/ws?jobID=exampleJobID");
console.log("Attempting Connection...");

socket.onopen = () => {
    console.log("Successfully Connected");
    setTimeout(() => {
        socket.send("scan")
    }
    , 1000)

};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};