using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddAuthorization();
builder.Services.AddRazorPages(options =>
{
    options.RootDirectory = "/Content";
    options.Conventions.AuthorizePage("/Admin/Settings");
});
var app = builder.Build();
app.UseRouting();
app.UseAuthorization();
app.MapRazorPages();
app.Run();
