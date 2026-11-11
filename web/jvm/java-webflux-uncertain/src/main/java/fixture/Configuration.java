package fixture;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.PropertySource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.server.session.*;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.reactive.UrlBasedCorsConfigurationSource;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;
import org.springframework.security.authorization.*;
import org.springframework.security.web.server.authorization.AuthorizationContext;
import org.thymeleaf.spring6.templateresolver.SpringResourceTemplateResolver;
import org.thymeleaf.templatemode.TemplateMode;
import java.util.function.UnaryOperator;
@org.springframework.context.annotation.Configuration
@PropertySource("classpath:cookie-latin1.properties")
@PropertySource(value = "classpath:cookie-settings.properties", encoding = "UTF-8")
@PropertySource("classpath:cookie-settings.xml")
@PropertySource(value = "classpath:cookie-settings.json", factory = JsonSettingsFactory.class)
public class Configuration {
  @Bean public SecurityWebFilterChain policyChain(ServerHttpSecurity http, ReactiveAuthorizationManager<AuthorizationContext> manager, UrlBasedCorsConfigurationSource source) {
    http.authorizeExchange(registry -> registry.anyExchange().access(manager)).cors(cors -> cors.configurationSource(source)).csrf(csrf -> csrf.disable()); return http.build();
  }
  @Bean public ReactiveAuthorizationManager<AuthorizationContext> manager() { return (authentication, context) -> authentication.map(auth -> new AuthorizationDecision(auth.isAuthenticated())); }
  @Bean(name = "webSessionManager") public WebSessionManager sessions(@Value("${web.cookie.secure}") boolean secure, @Value("${web.cookié.http-only}") boolean httpOnly) {
    CookieWebSessionIdResolver resolver = new CookieWebSessionIdResolver(); resolver.addCookieInitializer(cookie -> cookie.secure(secure).httpOnly(httpOnly));
    DefaultWebSessionManager manager = new DefaultWebSessionManager(); manager.setSessionIdResolver(resolver); return manager;
  }
  @Bean public UrlBasedCorsConfigurationSource origins(@Value("${web.origins}") String origins) {
    CorsConfiguration policy = new CorsConfiguration(); policy.setAllowCredentials(true); policy.setAllowedOriginPatterns(java.util.Arrays.asList(origins.split(",")));
    UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource(); source.registerCorsConfiguration("/**", policy); return source;
  }
  @Bean public UnaryOperator<String> renderer() { return value -> value.trim(); }
  @Bean public SpringResourceTemplateResolver xmlResolver() { SpringResourceTemplateResolver resolver = new SpringResourceTemplateResolver(); resolver.setPrefix("classpath:/templates/"); resolver.setSuffix(".xhtml"); resolver.setTemplateMode(TemplateMode.XML); resolver.setResolvablePatterns(java.util.Set.of("xml-*")); resolver.setOrder(0); return resolver; }
}
