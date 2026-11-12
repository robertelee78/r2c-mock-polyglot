using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc.RazorPages;

[Authorize]
public sealed class OrdersModel : PageModel
{
    public void OnGet()
    {
    }
}
