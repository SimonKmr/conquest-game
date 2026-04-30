import asyncio
import websockets
import json

class Settings:
    def __init__(self):
        self.uri = "ws://localhost:8080/player"
        self.id = 1
        self.password = "hello"

class Update:
    def __init__(self,id, x, y):
        self.id = id
        self.x = x
        self.y = y

class BotEngine:
    def reason(self, updates, settings):
        x = 126
        y = 77
        obj = {"Id": settings.id, "X": x, "Y": y}
        return obj


async def client():
    settings = Settings()
    engine = BotEngine()
    async with websockets.connect(
        settings.uri,
        open_timeout=None,    # handshake timeout (default: 10s)
        ping_timeout=None,    # pong wait timeout (default: 20s)
        close_timeout=None,   # graceful close timeout (default: 10s)
        ping_interval=None,   # disables keepalive pings entirely
    ) as websocket:
        # 1. Initialize — send first message
        await websocket.send(f"{settings.id} {settings.password}")
        print("Sent init message")

        # 2. Loop: receive then send
        while True:
            message = await websocket.recv()
            print(f"Received: {message}")

            updates = json.loads(message)
            reply = engine.reason(updates,settings)

            obj = json.dumps(reply)
            await websocket.send(obj)
            print(f"Sent: {obj}")



if __name__ == "__main__":
    asyncio.run(client())