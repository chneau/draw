window.onload = () => {
    console.log("Hello !");
    let canvas = document.getElementById('canvas');
    let ctx = canvas.getContext("2d");
    let socket = new WebSocket("ws://localhost:3000/ws");
}
