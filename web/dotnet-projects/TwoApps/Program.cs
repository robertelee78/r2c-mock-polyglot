using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

var adminBuilder = WebApplication.CreateBuilder(args);
adminBuilder.Services.AddRazorPages();
adminBuilder.Services.AddAuthorization();
var admin = adminBuilder.Build();
admin.UseRouting();
admin.UseAuthorization();
admin.MapRazorPages();

var siteBuilder = WebApplication.CreateBuilder(args);
siteBuilder.Services.AddRazorPages();
var site = siteBuilder.Build();
site.UseRouting();
site.MapRazorPages();

admin.Run();
