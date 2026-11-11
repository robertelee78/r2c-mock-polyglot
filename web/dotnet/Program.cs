using System.Collections.Generic;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.DependencyInjection;

var builder = WebApplication.CreateBuilder(args);
builder.Configuration.AddJsonFile("web-settings.json", optional: false);
builder.Services.AddControllersWithViews();
builder.Services.AddRazorPages();
builder.Services.AddDistributedMemoryCache();
builder.Services.AddSession(options => { options.Cookie.HttpOnly = true; });
builder.Services.AddSingleton<List<Record>>();
builder.Services.AddAntiforgery();
builder.Services.AddAuthorization();
builder.Services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
    .AddCookie("weak", options => { options.Cookie.HttpOnly = false; options.Cookie.SecurePolicy = CookieSecurePolicy.None; })
    .AddCookie("direct-control", options => { options.Cookie.HttpOnly = true; options.Cookie.SecurePolicy = CookieSecurePolicy.Always; })
    .AddCookie("selected-setting", options => { options.Cookie.HttpOnly = bool.Parse(builder.Configuration["Authentication:HttpOnly"]!); });
builder.Services.AddCors(options => {
    options.AddPolicy("arbitrary", policy => policy.SetIsOriginAllowed(origin => true).AllowCredentials());
    options.AddPolicy("finite", policy => policy.WithOrigins("https://trusted.example").AllowCredentials());
    options.AddPolicy("opaque", policy => policy.SetIsOriginAllowed(origin => OriginPolicy(origin)).AllowCredentials());
});
var app = builder.Build();
app.UseRouting();
app.UseSession();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();
app.MapRazorPages();

app.MapGet("/minimal/html/positive", (string value) => Results.Content(value, "text/html"));
app.MapGet("/minimal/html/control", (string value) => Results.Text(value));
app.MapGet("/minimal/html/uncertain", (string value) => Results.Content(Decorate(value), "text/html"));
app.MapGet("/minimal/redirect/positive", (string target) => Results.Redirect(target));
app.MapGet("/minimal/redirect/control", (string target) => Results.LocalRedirect(target));
app.MapGet("/minimal/redirect/uncertain", (string target) => Results.Redirect(Destination(target)));
app.MapPost("/minimal/csrf/positive", (HttpContext context) => {
    var ambient = context.Request.Cookies["session"];
    context.Session.SetString("state", ambient ?? "changed");
    return Results.Ok();
});
app.MapPost("/minimal/csrf/control", async (HttpContext context, Microsoft.AspNetCore.Antiforgery.IAntiforgery antiforgery) => {
    var ambient = context.Request.Cookies["session"];
    await antiforgery.ValidateRequestAsync(context);
    context.Session.SetString("state", ambient ?? "changed");
    return Results.Ok();
});
app.MapPost("/minimal/csrf/uncertain", (HttpContext context) => {
    var ambient = context.Request.Cookies["session"];
    Protect(context);
    context.Session.SetString("state", ambient ?? "changed");
    return Results.Ok();
});
app.MapPost("/minimal/authorization/positive", (HttpContext context) => {
    if (!context.User.IsInRole("approver")) Results.Forbid();
    context.Session.SetString("decision", "approved");
    return Results.Ok();
});
app.MapPost("/minimal/authorization/control", (HttpContext context) => {
    if (!context.User.IsInRole("approver")) return Results.Forbid();
    context.Session.SetString("decision", "approved");
    return Results.Ok();
});
app.MapPost("/minimal/authorization/uncertain", async (HttpContext context, IAuthorizationService authorization) => {
    var verdict = await authorization.AuthorizeAsync(context.User, null, "custom-policy");
    if (!verdict.Succeeded) return Results.Forbid();
    context.Session.SetString("decision", "approved");
    return Results.Ok();
});
app.MapPost("/minimal/binding/positive", ([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) => { records.Add(input); return Results.Ok(); });
app.MapPost("/minimal/binding/control", ([FromBody] Record input, [FromServices] List<Record> records) => { records.Add(new Record { DisplayName = input.DisplayName }); return Results.Ok(); });
app.MapPost("/minimal/binding/uncertain", ([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) => { ObservePolicy(input); records.Add(input); return Results.Ok(); });

var group = app.MapGroup("/group").RequireAuthorization();
group.MapGet("/html/positive", (string value) => Results.Content(value, "text/html"));
group.MapGet("/html/control", (string value) => Results.Text(value));
group.MapGet("/html/uncertain", (string value) => Results.Content(Decorate(value), "text/html"));
group.MapGet("/redirect/positive", (string target) => Results.Redirect(target));
group.MapGet("/redirect/control", (string target) => Results.LocalRedirect(target));
group.MapGet("/redirect/uncertain", (string target) => Results.Redirect(Destination(target)));
group.MapPost("/csrf/positive", (HttpContext context) => { var ambient = context.Request.Cookies["session"]; context.Session.SetString("state", ambient ?? "changed"); return Results.Ok(); });
group.MapPost("/csrf/control", async (HttpContext context, Microsoft.AspNetCore.Antiforgery.IAntiforgery antiforgery) => { var ambient = context.Request.Cookies["session"]; await antiforgery.ValidateRequestAsync(context); context.Session.SetString("state", ambient ?? "changed"); return Results.Ok(); });
group.MapPost("/csrf/uncertain", (HttpContext context) => { var ambient = context.Request.Cookies["session"]; Protect(context); context.Session.SetString("state", ambient ?? "changed"); return Results.Ok(); });
group.MapPost("/authorization/positive", (HttpContext context) => { if (!context.User.IsInRole("approver")) Results.Forbid(); context.Session.SetString("decision", "approved"); return Results.Ok(); });
group.MapPost("/authorization/control", (HttpContext context) => { if (!context.User.IsInRole("approver")) return Results.Forbid(); context.Session.SetString("decision", "approved"); return Results.Ok(); });
group.MapPost("/authorization/uncertain", async (HttpContext context, IAuthorizationService authorization) => { var verdict = await authorization.AuthorizeAsync(context.User, null, "custom-policy"); if (!verdict.Succeeded) return Results.Forbid(); context.Session.SetString("decision", "approved"); return Results.Ok(); });
group.MapPost("/binding/positive", ([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) => { records.Add(input); return Results.Ok(); });
group.MapPost("/binding/control", ([FromBody] Record input, [FromServices] List<Record> records) => { records.Add(new Record { DisplayName = input.DisplayName }); return Results.Ok(); });
group.MapPost("/binding/uncertain", ([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) => { ObservePolicy(input); records.Add(input); return Results.Ok(); });
app.Run();

static string Decorate(string value) => System.Net.WebUtility.HtmlDecode(value);
static string Destination(string value) => Uri.UnescapeDataString(value);
static bool OriginPolicy(string value) => Uri.TryCreate(value, UriKind.Absolute, out var uri) && uri.Host.EndsWith(".example", StringComparison.Ordinal);
static void Protect(HttpContext context) { context.Response.Headers.Append("X-Observed-Origin", context.Request.Headers.Origin); }
static void ObservePolicy(Record input) { input.DisplayName = input.DisplayName.Trim(); }
public sealed class Record { public string DisplayName { get; set; } = ""; public string Owner { get; set; } = "server"; }
