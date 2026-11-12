using Microsoft.AspNetCore.Builder;

public static class ConditionalEndpoints
{
    public static IResult Render(string value) => Results.Content(value, "text/html");
}
