using Microsoft.AspNetCore.Mvc.ModelBinding;
using System.Collections.Generic;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Antiforgery;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;

public sealed class MeasurementsModel : PageModel
{
    public string DisplayName { get; set; } = "";
    public IActionResult OnGet(string value) { DisplayName = value; return Page(); }
    public IActionResult OnGetHtmlPositive(string value) => Content(value, "text/html");
    public IActionResult OnGetHtmlControl(string value) => Content(value, "text/plain");
    public IActionResult OnGetHtmlUncertain(string value) => Content(System.Net.WebUtility.HtmlDecode(value), "text/html");
    public IActionResult OnGetRedirectPositive(string target) => Redirect(target);
    public IActionResult OnGetRedirectControl(string target) => LocalRedirect(target);
    public IActionResult OnGetRedirectUncertain(string target) => Redirect(Uri.UnescapeDataString(target));
    public IActionResult OnPostCsrfPositive() { var ambient = Request.Cookies["session"]; HttpContext.Session.SetString("state", ambient ?? "changed"); return Page(); }
    public async Task<IActionResult> OnPostCsrfControl([FromServices] IAntiforgery antiforgery) { var ambient = Request.Cookies["session"]; await antiforgery.ValidateRequestAsync(HttpContext); HttpContext.Session.SetString("state", ambient ?? "changed"); return Page(); }
    public IActionResult OnPostCsrfUncertain() { var ambient = Request.Cookies["session"]; Response.Headers.Append("X-Observed-Origin", Request.Headers.Origin); HttpContext.Session.SetString("state", ambient ?? "changed"); return Page(); }
    public IActionResult OnPostAuthorizationPositive() { if (!User.IsInRole("approver")) Forbid(); HttpContext.Session.SetString("decision", "approved"); return Page(); }
    public IActionResult OnPostAuthorizationControl() { if (!User.IsInRole("approver")) return Forbid(); HttpContext.Session.SetString("decision", "approved"); return Page(); }
    public async Task<IActionResult> OnPostAuthorizationUncertain([FromServices] IAuthorizationService authorization) { var verdict = await authorization.AuthorizeAsync(User, null, "custom-policy"); if (!verdict.Succeeded) return Forbid(); HttpContext.Session.SetString("decision", "approved"); return Page(); }
    public IActionResult OnPostBindingPositive([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) { records.Add(input); return Page(); }
    public IActionResult OnPostBindingControl([FromBody] Record input, [FromServices] List<Record> records) { records.Add(new Record { DisplayName = input.DisplayName }); return Page(); }
    public IActionResult OnPostBindingUncertain([FromBody, BindNever] Record input, [FromServices] List<Record> records) { input.DisplayName = input.DisplayName.Trim(); records.Add(input); return Page(); }
}
