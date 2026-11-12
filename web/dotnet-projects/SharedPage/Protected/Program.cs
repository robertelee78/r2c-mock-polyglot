using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddRazorPages();
builder.Services.AddAuthorization();
var app = builder.Build();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
