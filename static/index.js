let kelly = ["#e6194b", "#3cb44b", "#ffe119", "#0082c8", "#f58231", "#911eb4", "#46f0f0", "#f032e6", "#d2f53c", "#fabebe", "#008080", "#e6beff", "#aa6e28", "#fffac8", "#800000", "#aaffc3", "#808000", "#ffd8b1", "#000080", "#808080", "#FFFFFF", "#000000"];

window.onload = () => {
    let canvas = document.getElementById('canvas');
    let ctx = canvas.getContext("2d");
    ctx.canvas.width = window.innerWidth;
    ctx.canvas.height = window.innerHeight;
    let socket = new WebSocket("ws://localhost:3000/ws");
    let isDrawing = false;
    let color = 0;
    socket.onmessage = (event) => {
        let msg = JSON.parse(event.data);
        ctx.fillStyle = kelly[msg.c];
        ctx.fillRect(msg.x, msg.y, 15, 15);
    };
    canvas.onmousemove = (event) => {
        if (event.which != 0) {
            socket.send(JSON.stringify({ x: event.pageX, y: event.pageY, c: color % kelly.length }));
        }
    };
    canvas.oncontextmenu = (event) => {
        event.preventDefault();
    };
    canvas.onmousewheel = (event) => {
        if (++color <= 0) {
            color = 0;
        }
        event.preventDefault();
    }

}
