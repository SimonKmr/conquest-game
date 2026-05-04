namespace client_cs
{
    public class Update
    {
        public int Id { get; set; }
        public int X { get; set; }
        public int Y { get; set; }

        public override string ToString()
        {
            return $"{Id} - {X} ; {Y}";
        }
    }
}
