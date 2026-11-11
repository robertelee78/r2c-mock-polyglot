package fixture;
import org.springframework.web.bind.annotation.*;
import org.springframework.http.*;
import org.springframework.web.util.HtmlUtils;
import org.springframework.web.servlet.view.RedirectView;
import org.springframework.web.server.WebSession;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.security.web.csrf.CsrfToken;
import org.springframework.security.web.csrf.DefaultCsrfToken;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.beans.BeanUtils;
import jakarta.servlet.http.HttpSession;
import java.util.function.UnaryOperator;
@RestController
public class Cases {
  private final CrudRepository<Account, Long> repository;
  private final UnaryOperator<String> renderer;
  public Cases(CrudRepository<Account, Long> repository, UnaryOperator<String> renderer) { this.repository = repository; this.renderer = renderer; }
  @GetMapping("/html")
  public ResponseEntity<String> html(@RequestParam String input) { return ResponseEntity.ok().contentType(MediaType.TEXT_HTML).body(input); }
  @GetMapping("/redirect")
  public RedirectView redirect(@RequestParam String target) { return new RedirectView(target); }
  @GetMapping("/token")
  public String token(HttpSession session) { String token = java.util.UUID.randomUUID().toString(); session.setAttribute("formToken", token); return token; }
  @PostMapping("/state")
  public String change(HttpSession session, @RequestHeader("X-CSRF-TOKEN") String provided, @RequestParam String value) {
    CsrfToken token = new DefaultCsrfToken("X-CSRF-TOKEN", "_csrf", (String) session.getAttribute("formToken"));
    
    session.setAttribute("displayName", value); return "changed";
  }
  @GetMapping("/policy") @PreAuthorize("hasRole('OPERATOR')")
  public String policy() { return "policy protected"; }
  @PostMapping("/account")
  public Account account(@RequestBody InputAccount input, @AuthenticationPrincipal UserDetails principal) { Account target = new Account(); target.setOwner(principal.getUsername()); BeanUtils.copyProperties(input, target); return repository.save(target); }
  private boolean accepted(String value) { return value != null && value.length() > 20; }
}
