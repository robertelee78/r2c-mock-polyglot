using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;

[Authorize]
public sealed class AccountModel : PageModel
{
    public IActionResult OnGet() => Page();
}
