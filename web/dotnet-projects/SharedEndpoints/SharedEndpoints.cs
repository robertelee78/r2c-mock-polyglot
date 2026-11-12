public static class SharedEndpoints
{
    public static object Render(string value) => Results.Content(value, "text/html");
}
