using Microsoft.AspNetCore.Mvc.ModelBinding;
using System.Collections.Generic;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Antiforgery;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;

[Route("mvc")]
public sealed class MeasurementsController : Controller
{
    [HttpGet("html/positive")]
    public IActionResult HtmlPositive(string value) => Content(value, "text/html");
    [HttpGet("html/control")]
    public IActionResult HtmlControl(string value) => Content(value, "text/plain");
    [HttpGet("html/uncertain")]
    public IActionResult HtmlUncertain(string value) => View("Raw", value);
    [HttpGet("razor/control")]
    public IActionResult RazorControl(string value) => View("Encoded", value);
    [HttpGet("redirect/positive")]
    public IActionResult RedirectPositive(string target) => Redirect(target);
    [HttpGet("redirect/control")]
    public IActionResult RedirectControl(string target) => LocalRedirect(target);
    [HttpGet("redirect/uncertain")]
    public IActionResult RedirectUncertain(string target) => Redirect(Uri.UnescapeDataString(target));
    [HttpPost("csrf/positive"), IgnoreAntiforgeryToken]
    public IActionResult CsrfPositive() { var ambient = Request.Cookies["session"]; HttpContext.Session.SetString("state", ambient ?? "changed"); return Ok(); }
    [HttpPost("csrf/control"), ValidateAntiForgeryToken]
    public IActionResult CsrfControl() { var ambient = Request.Cookies["session"]; HttpContext.Session.SetString("state", ambient ?? "changed"); return Ok(); }
    [HttpPost("csrf/uncertain")]
    public IActionResult CsrfUncertain() { var ambient = Request.Cookies["session"]; Response.Headers.Append("X-Observed-Origin", Request.Headers.Origin); HttpContext.Session.SetString("state", ambient ?? "changed"); return Ok(); }
    [HttpPost("authorization/positive")]
    public IActionResult AuthorizationPositive() { if (!User.IsInRole("approver")) Forbid(); HttpContext.Session.SetString("decision", "approved"); return Ok(); }
    [HttpPost("authorization/control")]
    public IActionResult AuthorizationControl() { if (!User.IsInRole("approver")) return Forbid(); HttpContext.Session.SetString("decision", "approved"); return Ok(); }
    [HttpPost("authorization/uncertain")]
    public async Task<IActionResult> AuthorizationUncertain([FromServices] IAuthorizationService authorization) { var verdict = await authorization.AuthorizeAsync(User, null, "custom-policy"); if (!verdict.Succeeded) return Forbid(); HttpContext.Session.SetString("decision", "approved"); return Ok(); }
    [HttpPost("binding/positive")]
    public IActionResult BindingPositive([FromBody, Bind("DisplayName")] Record input, [FromServices] List<Record> records) { records.Add(input); return Ok(); }
    [HttpPost("binding/control")]
    public IActionResult BindingControl([FromBody] Record input, [FromServices] List<Record> records) { records.Add(new Record { DisplayName = input.DisplayName }); return Ok(); }
    [HttpPost("binding/uncertain")]
    public IActionResult BindingUncertain([FromBody, BindNever] Record input, [FromServices] List<Record> records) { input.DisplayName = input.DisplayName.Trim(); records.Add(input); return Ok(); }
}
