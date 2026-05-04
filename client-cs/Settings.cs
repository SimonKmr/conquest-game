namespace client_cs
{
    public class Settings
    {
        public int Id { get; set; }
        public string Uri { get; set; }
        public string Password { get; set; }

        public Settings(int id, string uri, string password)
        {
            Id = id;
            Uri = uri;
            Password = password;
        }
    }
}
