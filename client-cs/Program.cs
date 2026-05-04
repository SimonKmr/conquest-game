using WebSocketSharp;
using Newtonsoft.Json;

namespace client_cs
{
    internal class Program
    {
        static void Main(string[] args)
        {
            Settings settings = new(2, "ws://localhost:8080/player", "world");
            BotEngine engine = new();
            using var ws = new WebSocket(settings.Uri) ;

            ws.WaitTime = TimeSpan.FromMinutes(60);
            ws.OnMessage += (sender, e) =>
            {
                Console.WriteLine(e.Data);
                var updates = JsonConvert.DeserializeObject<Update[]>(e.Data);
                var move = engine.Reason(settings, updates);
                var moveJson = JsonConvert.SerializeObject(move);
                Console.WriteLine(moveJson);
                ws.Send(moveJson);
                    
            };
            ws.Connect();
            ws.Send($"{settings.Id} {settings.Password}");
            while (true) ;

        }
    }
}
