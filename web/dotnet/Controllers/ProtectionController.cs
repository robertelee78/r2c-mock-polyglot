using Microsoft.AspNetCore.Antiforgery;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;

[ValidateAntiForgeryToken]
[Route("protection")]
public sealed class ProtectionController : Controller
{
    [HttpPost("class-antiforgery")]
    public IActionResult ClassAntiforgery() { var ambient = Request.Cookies["session"]; HttpContext.Session.SetString("state", ambient ?? "changed"); return Ok(); }
    [HttpGet("encoded-local")]
    public IActionResult EncodedLocal(string value) { var encoded = System.Net.WebUtility.HtmlEncode(value); return Content(encoded, "text/html"); }
    [HttpGet("literal-local")]
    public IActionResult LiteralLocal() { var home = "/home"; return Redirect(home); }
    [HttpGet("guarded")]
    public IActionResult Guarded(string target, bool verbose) { if (verbose) { target = target.Trim(); } if (!Url.IsLocalUrl(target)) return BadRequest(); return Redirect(target); }
}

[AutoValidateAntiforgeryToken]
[Route("auto")]
public sealed class AutoProtectionController : Controller
{
    [HttpPost("state")]
    public IActionResult State() { var ambient = Request.Cookies["session"]; HttpContext.Session.SetString("state", ambient ?? "changed"); return Ok(); }
}
