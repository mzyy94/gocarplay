const video = document.querySelector("video");
const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: "stun:stun.l.google.com:19302",
    },
  ],
});

pc.ontrack = (event) => {
  if (video.srcObject == null) {
    video.srcObject = event.streams[0];
  } else {
    video.srcObject.addTrack(event.track);
  }
};

const startData = pc.createDataChannel("start");
startData.onopen = () => startData.send(
  JSON.stringify({
    width: (innerWidth * devicePixelRatio) | 0,
    height: (innerHeight * devicePixelRatio) | 0,
  })
);

pc.oniceconnectionstatechange = () => {
  console.log("connection:", pc.iceConnectionState);
};

pc.onicecandidate = (event) => {
  if (event.candidate == null) {
    fetch("/connect", {
      method: "POST",
      body: JSON.stringify(pc.localDescription),
    })
      .then((res) => Promise.all([res.json(), res.ok]))
      .then(([answer, ok]) => {
        if (!ok) {
          return Promise.reject(answer);
        }
        try {
          pc.setRemoteDescription(new RTCSessionDescription(answer));
        } catch (e) {
          Promise.reject(e);
        }
      })
      .catch(console.error);
  }
};

pc.addTransceiver("video", { direction: "recvonly" });

const audioCtx = new (window.AudioContext || window.webkitAudioContext)();

pc.ondatachannel = ({ channel: dc }) => {
  if (dc.label == "audio") {
    dc.onmessage = (e) => {
      const dv = new DataView(e.data.slice(0, 4));
      const data = new Float32Array(new Int16Array(e.data.slice(4))).map(
        (d) => d / 32768
      );
      const sampleRate = dv.getUint16(0, true);
      const channels = dv.getUint16(2, true);
      const audioBuffer = audioCtx.createBuffer(
        channels,
        data.length / channels,
        sampleRate
      );

      for (let ch = 0; ch < channels; ++ch) {
        audioBuffer
          .getChannelData(ch)
          .set(data.filter((_, i) => i % channels == ch));
      }

      const src = audioCtx.createBufferSource();
      src.buffer = audioBuffer;
      src.connect(audioCtx.destination);
      src.start();
    };
  }
};

const touchData = pc.createDataChannel("touch");

let pointerdown = false;
const sendTouchEvent = ({ type, offsetX, offsetY }) => {
  let action = 16;
  if (type == "pointerdown") {
    action = 14;
    pointerdown = true;
  } else if (pointerdown) {
    switch (type) {
      case "pointermove":
        action = 15;
        break;
      case "pointerup":
      case "pointercancel":
      case "pointerout":
        pointerdown = false;
        action = 16;
        break;
    }
  } else {
    return;
  }
  const data = {
    x: (offsetX * devicePixelRatio) | 0,
    y: (offsetY * devicePixelRatio) | 0,
    action,
  };
  touchData.send(JSON.stringify(data));
};

video.addEventListener("pointerdown", sendTouchEvent);
video.addEventListener("pointermove", sendTouchEvent);
video.addEventListener("pointerup", sendTouchEvent);
video.addEventListener("pointercancel", sendTouchEvent);
video.addEventListener("pointerout", sendTouchEvent);

pc.createOffer()
  .then((d) => pc.setLocalDescription(d))
  .catch(console.error);
