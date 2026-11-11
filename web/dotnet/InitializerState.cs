using System;
using Microsoft.AspNetCore.Html;

public static class InitializerState
{
    public static readonly HtmlString StartupHtml = new HtmlString(ReadInput());
    public static HtmlString PropertyHtml { get; } = new HtmlString(ReadInput());
    public const string Label = "constant-only data";
    private static string ReadInput() => Environment.GetEnvironmentVariable("O84_HEADER") ?? "";
}
