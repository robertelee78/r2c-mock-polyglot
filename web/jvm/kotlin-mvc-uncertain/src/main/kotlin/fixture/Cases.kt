package fixture
import org.springframework.web.bind.annotation.*
import org.springframework.http.*
import org.springframework.web.util.HtmlUtils
import org.springframework.web.servlet.view.RedirectView
import org.springframework.web.server.WebSession
import org.springframework.web.server.ResponseStatusException
import org.springframework.security.web.csrf.CsrfToken
import org.springframework.security.web.csrf.DefaultCsrfToken
import org.springframework.security.access.prepost.PreAuthorize
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity
import org.springframework.data.repository.CrudRepository
import org.springframework.beans.BeanUtils
import jakarta.servlet.http.HttpSession
import java.util.function.UnaryOperator
@RestController
open class Cases(private val repository: CrudRepository<Account, Long>, private val renderer: UnaryOperator<String>) {
  @GetMapping("/html")
  open fun html(@RequestParam input: String): ResponseEntity<String> { return ResponseEntity.ok().contentType(MediaType.TEXT_HTML).body(renderer.apply(input)) }
  @GetMapping("/redirect")
  open fun redirect(@RequestParam target: String): RedirectView { return RedirectView(renderer.apply(target)) }
  @GetMapping("/token")
  open fun token(session: HttpSession): String { val token: String = java.util.UUID.randomUUID().toString(); session.setAttribute("formToken", token); return token }
  @PostMapping("/state")
  open fun change(session: HttpSession, @RequestHeader("X-CSRF-TOKEN") provided: String, @RequestParam value: String): String {
    val token: CsrfToken = DefaultCsrfToken("X-CSRF-TOKEN", "_csrf", session.getAttribute("formToken") as String)
    if (!accepted(provided)) { throw ResponseStatusException(HttpStatus.FORBIDDEN) }
    session.setAttribute("displayName", value); return "changed"
  }
  @GetMapping("/policy") @PreAuthorize("hasRole('OPERATOR')")
  open fun policy(): String { return "policy protected" }
  @PostMapping("/account")
  open fun account(@ModelAttribute input: Account, @AuthenticationPrincipal principal: UserDetails): Account { return repository.save(input) }
  private fun accepted(value: String): Boolean { return value.length > 20 }
}
