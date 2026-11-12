using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddCors(options =>
{
    options.AddPolicy("listed", policy => policy.SetIsOriginAllowed(origin => origin == "https://a.example" || origin == "https://b.example").AllowCredentials());
});
var app = builder.Build();
app.UseCors("listed");
app.Run();
